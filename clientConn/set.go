package clientConn

import (
	"MaybeDB/server/cluster"
	"MaybeDB/server/database"
	"MaybeDB/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

//插入数据
func Set(c *gin.Context) {

	isCluster := viper.GetInt("server.isCluster")
	isMaster := viper.GetInt("server.isMaster")

	//集群模式下，只有主节点可以写入/删除数据，从节点只负责读取数据
	if isMaster == 0 && isCluster == 1 {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "only the master node can write data",
		})
		return
	}

	key, _ := c.GetPostForm("key")
	value, _ := c.GetPostForm("value")
	valueType, _ := c.GetPostForm("valueType")
	expireTime, _ := c.GetPostForm("expireTime")

	valueTypeInt, err := strconv.Atoi(valueType)
	if err != nil {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  err,
		})
		return
	}

	expireTimeInt64, err := strconv.ParseInt(expireTime, 10, 64)
	if err != nil {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  err,
		})
		return
	}

	var expTime int64

	if expireTimeInt64 == 0 {
		expTime = 0
	} else {
		expTime = time.Now().Unix() + expireTimeInt64
	}

	var content interface{}

	switch valueTypeInt {

	//string类型的值
	case 1:
		content = value
		break

	//int64类型的值
	case 2:
		//string转int
		content, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			c.JSON(-1, gin.H{
				"code": -1,
				"msg":  err,
			})
			return
		}
		break

	//map类型的值（value为json字符串，例：`{"id": "1", "text": "hello"}`）
	case 3:
		//json转map
		content, err = utils.JsonToMap(value)
		if err != nil {
			c.JSON(-1, gin.H{
				"code": -1,
				"msg":  err,
			})
			return
		}
		break

	//array类型的值（value为json字符串，例：`[1, 2, 3, 4, 5]`）
	case 4:
		//将json字符串解析为数组
		err = json.Unmarshal([]byte(value), &content)
		if err != nil {
			c.JSON(-1, gin.H{
				"code": -1,
				"msg":  err,
			})
			return
		}
		break

	//无效类型
	default:
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "invalid type",
		})
		return
	}

	//生成数据模板
	data := database.Data{
		Content:    content,
		ExpireTime: expTime,
	}

	//插入数据
	database.DataMap.Store(key, data)

	// isOtherMaster：判断该数据是否是其他主节点发来的
	isOtherMaster := c.GetHeader("isOtherMaster")
	//如果是集群模式，且该消息不是从其他主节点发来的
	if isCluster == 1 && isOtherMaster != "1" {
		//将该新增数据同步到其他主节点
		masterSetSync(key, value, valueType, expireTime)
	}

	c.JSON(0, gin.H{
		"code": 0,
	})
}

//将新增数据同步到所有主节点
func masterSetSync(key string, value string, valueType string, expireTime string) {

	//获取该主节点的地址
	ip := viper.GetString("server.ip")
	port := uint64(viper.GetInt("server.port"))

	//获取健康的主节点实例列表
	instances, _ := cluster.NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: "maybe-db-master",
		GroupName:   "MAYBE_DB_GROUP",             // 默认值DEFAULT_GROUP
		Clusters:    []string{"MAYBE_DB_CLUSTER"}, // 默认值DEFAULT
		HealthyOnly: true,
	})

	size := len(instances)
	if size < 2 {
		return
	}

	//设置请求头
	header := make(map[string]string, 2)
	//标记，表明该请求是主节点集群同步数据，不是客户端发来的数据
	header["isOtherMaster"] = "1"
	//访问密钥
	header["secretKey"] = viper.GetString("db.secretKey")

	//设置请求数据
	data := make(map[string]string, 4)
	data["key"] = key
	data["value"] = value
	data["valueType"] = valueType
	data["expireTime"] = expireTime

	//控制通道，用于确保同步工作完成后，再return
	sendChan := make(chan int, size)

	//同步向所有主节点发送新增数据
	for _, instance := range instances {

		//开启协程
		go func(instance model.Instance) {
			if instance.Ip == ip && instance.Port == port {
				//如果遍历到该主节点自己，则跳过
				sendChan <- 1
			} else {
				//向其他主节点发送数据
				url := fmt.Sprintf("%s%s%s%s%s", "http://", instance.Ip, ":", strconv.Itoa(int(instance.Port)), "/Client/Set")
				_, err := utils.PostForm(url, header, data)
				if err != nil {
					database.Loger.Println(err)
				}
				sendChan <- 1
			}
		}(instance)
	}

	//所有请求都发送完毕后，再return
	for range instances {
		<-sendChan
	}
}
