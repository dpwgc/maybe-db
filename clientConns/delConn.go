package clientConns

import (
	"MaybeDB/cluster"
	"MaybeDB/servers"
	"MaybeDB/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"strconv"
)

//删除数据
func Del(c *gin.Context) {

	isCluster := viper.GetInt("server.isCluster")
	isMaster := viper.GetInt("server.isMaster")

	//集群模式下，只有主节点可以写入/删除数据，从节点只负责读取数据
	if isMaster == 0 && isCluster == 1 {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "only the master node can delete data",
		})
		return
	}

	key, _ := c.GetPostForm("key")

	servers.DataMap.Delete(key)

	// isOtherMaster：判断该数据是否是其他主节点发来的
	isOtherMaster := c.GetHeader("isOtherMaster")
	//如果是集群模式，且该消息不是从其他主节点发来的
	if isCluster == 1 && isOtherMaster != "1" {
		//将删除操作同步到其他主节点
		masterDelSync(key)
	}

	c.JSON(0, gin.H{
		"code": 0,
	})
}

//将删除操作同步到所有主节点
func masterDelSync(key string) {

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

	//设置请求头
	header := make(map[string]string, 2)
	//标记，表明该请求是主节点集群同步数据，不是客户端发来的数据
	header["isOtherMaster"] = "1"
	//访问密钥
	header["secretKey"] = viper.GetString("db.secretKey")

	//设置请求数据
	data := make(map[string]string, 4)
	data["key"] = key

	//向所有主节点发送新增数据
	for _, instance := range instances {

		if instance.Ip == ip && instance.Port == port {
			//如果遍历到该主节点自己，则跳过
			continue
		}
		url := fmt.Sprintf("%s%s%s%s%s", "http://", instance.Ip, ":", strconv.Itoa(int(instance.Port)), "/Client/Del")
		resStr, err := utils.PostForm(url, header, data)

		res := make(map[string]interface{})
		err = json.Unmarshal([]byte(resStr), &res)
		if err != nil {
			servers.Loger.Println(err)
		}
	}
}
