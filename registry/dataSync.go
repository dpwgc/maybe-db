package registry

import (
	"MaybeDB/servers"
	"MaybeDB/utils"
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"time"
)

/**
 * 数据同步
 */

//数据库节点间数据同步
func SyncInit() {
	syncTime := viper.GetInt("db.syncTime")
	go func() {
		for {
			copyDataMap()
			servers.ByteCopyMap, _ = json.Marshal(servers.CopyMap)
			servers.JsonCopyMap = string(servers.ByteCopyMap)
			updateMetadata()
			pullMetadata()
			time.Sleep(time.Second * time.Duration(syncTime))
		}
	}()
}

//主节点更新nacos元数据
func updateMetadata() {
	namingClient.UpdateInstance(vo.UpdateInstanceParam{
		Ip:          viper.GetString("server.ip"),
		Port:        uint64(viper.GetInt("server.port")),
		ServiceName: "maybe-db",
		Weight:      10,
		Enable:      true,
		Ephemeral:   true,
		Metadata:    map[string]string{"DataMap": servers.JsonCopyMap},
		ClusterName: "MAYBE_DB",
		GroupName:   "MAYBE_DB_GROUP",
	})
}

//从节点拉取主节点的元数据，并更新本地DataMap
func pullMetadata() {
	//获取nacos存在服务的信息
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "maybe-db",
		GroupName:   "MAYBE_DB_GROUP",
		Clusters:    []string{"MAYBE_DB"},
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(utils.JsonToMap(instance.Metadata["DataMap"]))
	}
}

//复制数据集合
func copyDataMap() {
	servers.DataMap.Range(func(key, value interface{}) bool {
		servers.CopyMap[key.(string)] = value
		return true
	})
}
