package cluster

import (
	"MaybeDB/server/database"
	"github.com/spf13/viper"
	"time"
)

/**
 * 主从节点数据同步
 */

//主从节点之间的数据同步
func SyncInit() {

	//是否以集群方式部署
	isCluster := viper.GetInt("server.isCluster")
	if isCluster == 0 {
		return
	}

	//加载配置信息
	syncTime := viper.GetInt("server.syncTime")
	isMaster := viper.GetInt("server.isMaster")

	//开启主从同步协程
	go func() {
		//如果该节点是主节点
		if isMaster == 1 {
			database.Loger.Println("Master starts synchronization")
			for {
				//定时复制本地DataMap
				copyDataMap()
				time.Sleep(time.Second * time.Duration(syncTime))
			}
		}

		//如果该节点是从节点
		if isMaster == 0 {
			for {
				database.Loger.Println("Slave starts synchronization")
				//从节点定时请求获取主节点的数据，并同步更新本地DataMap
				syncWithMaster()
				time.Sleep(time.Second * time.Duration(syncTime))
			}
		}
		return
	}()
}
