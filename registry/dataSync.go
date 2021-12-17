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

//数据库主从节点之间的数据同步进程
func SyncInit() {
	syncTime := viper.GetInt("db.syncTime")
	isMaster := viper.GetInt("db.isMaster")
	go func() {
		//如果该节点是主节点
		if isMaster == 1 {
			for {
				//复制本地DataMap
				copyDataMap()
				//更新nacos元数据
				updateMetadata()
				time.Sleep(time.Second * time.Duration(syncTime))
			}
		}

		//如果该节点是从节点
		if isMaster == 0 {
			for {
				//拉取主节点的nacos元数据信息，将元数据写入本地DataMap
				pullMetadata()
				time.Sleep(time.Second * time.Duration(syncTime))
			}
		}
		return
	}()
}

//主节点操作
//主节点更新nacos元数据(将主节点的DataMap拷贝到nacos元数据上)
func updateMetadata() {
	namingClient.UpdateInstance(vo.UpdateInstanceParam{
		Ip:          viper.GetString("server.ip"),
		Port:        uint64(viper.GetInt("server.port")),
		ServiceName: "maybe-db",
		Weight:      10,
		Enable:      true,
		Ephemeral:   true,
		Metadata:    map[string]string{"DataMap": servers.JsonCopyMap}, //将主节点DataMap以Json字符串形式存入nacos元数据
		ClusterName: "MAYBE_DB",
		GroupName:   "MAYBE_DB_GROUP",
	})
}

//从节点操作
//从节点拉取主节点的元数据，并更新本地DataMap
func pullMetadata() {
	//获取nacos上数据库主节点的元数据
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "maybe-db",
		GroupName:   "MAYBE_DB_GROUP",
		Clusters:    []string{"MAYBE_DB"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//解析元数据
	masterMap, err := utils.JsonToData(instance.Metadata["DataMap"])
	if err != nil {
		fmt.Println(err)
		return
	}

	//将主节点map内的数据循环写入从节点map
	for key, value := range masterMap {
		servers.DataMap.Store(key, value)
	}
}

//复制数据集合
func copyDataMap() {
	servers.DataMap.Range(func(key, value interface{}) bool {
		servers.CopyMap[key.(string)] = value
		return true
	})
	//将CopyMap转为字节数组类型ByteCopyMap
	servers.ByteCopyMap, _ = json.Marshal(servers.CopyMap)
	//将ByteCopyMap转为Json字符串类型
	servers.JsonCopyMap = string(servers.ByteCopyMap)
}
