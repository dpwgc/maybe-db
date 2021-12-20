package cluster

import (
	"MaybeDB/server/database"
	"MaybeDB/utils"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"strconv"
)

/*
 * 主从同步，从节点操作
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
		database.Loger.Println(err)
		return
	}

	//设置请求头
	header := make(map[string]string, 1)
	//访问密钥
	header["secretKey"] = viper.GetString("db.secretKey")

	//向该主节点请求数据
	url := fmt.Sprintf("%s%s%s%s%s", "http://", instance.Ip, ":", strconv.Itoa(int(instance.Port)), "/Sync/GetMasterData")
	res, err := utils.Get(url, header)
	if err != nil {
		database.Loger.Println(err)
		return
	}

	//解析数据到masterMap集合
	masterMap, err := utils.JsonToData(res)
	if err != nil {
		database.Loger.Println(err)
		return
	}

	//将主节点map内的数据（masterMap）循环写入从节点map（DataMap）
	for key, value := range masterMap {
		database.DataMap.Store(key, value)
	}
}
