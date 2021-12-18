package cluster

import (
	"MaybeDB/servers"
	"MaybeDB/utils"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

/*
 * 主节点数据恢复
 */

//主节点重启后进行数据恢复（集群上必须要有其他健康的主节点实例，否则无法恢复）
func RecoveryInit() {

	//获取一个健康的主节点实例，获取主节点上的Nacos元数据
	instance, err := NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "maybe-db-master",
		GroupName:   "MAYBE_DB_GROUP",
		Clusters:    []string{"MAYBE_DB_CLUSTER"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//解析元数据到masterMap集合
	masterMap, err := utils.JsonToData(instance.Metadata["DataMap"])
	if err != nil {
		fmt.Println(err)
		return
	}

	//将主节点map内的数据（masterMap）循环写入从节点map（DataMap）
	for key, value := range masterMap {
		servers.DataMap.Store(key, value)
	}
}
