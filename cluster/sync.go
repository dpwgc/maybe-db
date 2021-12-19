package cluster

import (
	"github.com/spf13/viper"
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
				//定时复制本地DataMap
				copyDataMap()
				time.Sleep(time.Second * time.Duration(syncTime))
			}
		}

		//如果该节点是从节点
		if isMaster == 0 {
			for {
				//从节点定时请求获取主节点的数据，并同步更新本地DataMap
				syncWithMaster()
				time.Sleep(time.Second * time.Duration(syncTime))
			}
		}
		return
	}()
}
