package cluster

import (
	"MaybeDB/servers"
	"MaybeDB/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

/**
 * 主从节点数据同步
 */

//主从节点之间的数据同步
func SyncInit() {

	//加载配置信息
	syncTime := viper.GetInt("server.syncTime")
	isMaster := viper.GetInt("server.isMaster")

	//开启主从同步协程
	go func() {
		//如果该节点是主节点
		if isMaster == 1 {
			for {
				//复制本地DataMap
				copyDataMap()
				time.Sleep(time.Second * time.Duration(syncTime))
			}
		}

		//如果该节点是从节点
		if isMaster == 0 {
			for {
				//从节点请求获取主节点的数据，并同步更新本地DataMap
				syncWithMaster()
				time.Sleep(time.Second * time.Duration(syncTime))
			}
		}
		return
	}()
}

/*
 *主节点操作
 */

//该主节点的DataMap数据获取接口，用于提供给从节点
func GetMasterData(c *gin.Context) {

	//以Json字符串形式返回主节点的全部数据
	res := servers.SyncCopyJson

	c.JSON(0, gin.H{
		"code": 0,
		"data": res,
	})
}

//复制该主节点的本地数据
func copyDataMap() {
	servers.DataMap.Range(func(key, value interface{}) bool {
		servers.SyncCopyMap[key.(string)] = value
		return true
	})
	//将SyncCopyMap转为字节数组类型SyncCopyByte
	servers.SyncCopyByte, _ = json.Marshal(servers.SyncCopyMap)
	//将SyncCopyByte转为Json字符串类型
	servers.SyncCopyJson = string(servers.SyncCopyByte)
}

/*
 *从节点操作
 */

//从节点请求获取主节点的数据，并同步更新本地DataMap
func syncWithMaster() {

	//获取一个健康的主节点实例
	instance, err := NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "maybe-db-master",
		GroupName:   "MAYBE_DB_GROUP",
		Clusters:    []string{"MAYBE_DB_CLUSTER"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//设置请求头
	header := make(map[string]string, 1)
	//访问密钥
	header["secretKey"] = viper.GetString("db.secretKey")

	//向该主节点请求数据
	url := fmt.Sprintf("%s%s%s%s%s", "http://", instance.Ip, ":", strconv.Itoa(int(instance.Port)), "/Sync/GetMasterData")
	resStr, err := utils.Get(url, header)

	var res string
	err = json.Unmarshal([]byte(resStr), &res)
	if err != nil {
		fmt.Println(err)
	}

	//解析数据到masterMap集合
	masterMap, err := utils.JsonToData(res)
	if err != nil {
		fmt.Println(err)
		return
	}

	//将主节点map内的数据（masterMap）循环写入从节点map（DataMap）
	for key, value := range masterMap {
		servers.DataMap.Store(key, value)
	}
}
