package servers

import (
	"MaybeDB/utils"
	"encoding/json"
	"github.com/spf13/viper"
	"time"
)

/*
 * 持久化数据
 */

func InitPers() {

	//两次持久化操作的间隔时间
	persistentTime := viper.GetInt("db.persistentTime")

	go func() {
		//复制本地数据
		copyDataMap()
		//持久化
		utils.Write()
		time.Sleep(time.Second * time.Duration(persistentTime))
	}()
}

//复制该主节点的本地数据
func copyDataMap() {
	DataMap.Range(func(key, value interface{}) bool {
		PersCopyMap[key.(string)] = value
		return true
	})
	//将PersCopyMap转为字节数组类型PersCopyByte
	PersCopyByte, _ = json.Marshal(PersCopyMap)
	//将PersCopyByte转为Json字符串类型
	PersCopyJson = string(PersCopyByte)
}
