package recovery

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

	//两次持久化操作的间隔时间
	persistentTime := viper.GetInt("db.persistentTime")

	go func() {
		for {
			//复制本地数据
			copyDataMap()
			//持久化
			Write()
			time.Sleep(time.Second * time.Duration(persistentTime))
		}
	}()
}

//复制该主节点的数据
func copyDataMap() {
	servers.DataMap.Range(func(key, value interface{}) bool {
		servers.PersCopyMap[key.(string)] = value
		return true
	})
	//将PersCopyMap转为字节数组类型PersCopyByte
	servers.PersCopyByte, _ = json.Marshal(servers.PersCopyMap)
	//将PersCopyByte转为Json字符串类型
	servers.PersCopyJson = string(servers.PersCopyByte)
}
