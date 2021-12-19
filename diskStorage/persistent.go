package diskStorage

import (
	"MaybeDB/servers"
	"encoding/json"
	"github.com/spf13/viper"
	"time"
)

/*
 * 持久化数据到硬盘
 */

func PersInit() {

	//是否以集群方式部署
	isCluster := viper.GetInt("server.isCluster")
	//是否是主节点
	isMaster := viper.GetInt("server.isMaster")

	//集群模式下，从节点会自动同步主节点数据，无需进行数据持久化
	if isCluster == 1 && isMaster == 0 {
		return
	}

	//两次持久化操作的间隔时间
	persistentTime := viper.GetInt("db.persistentTime")

	go func() {
		servers.Loger.Println("Start persistence")
		for {
			//复制本地数据
			copyDataMap()
			//持久化写入
			Write()
			time.Sleep(time.Second * time.Duration(persistentTime))
		}
	}()
}

//复制该主节点的数据
func copyDataMap() {

	copyMap := make(map[string]interface{})

	servers.DataMap.Range(func(key, value interface{}) bool {
		copyMap[key.(string)] = value
		return true
	})
	//将PersCopyMap转为字节数组类型PersCopyByte
	servers.PersCopyByte, _ = json.Marshal(copyMap)
}
