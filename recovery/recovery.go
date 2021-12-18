package recovery

import (
	"MaybeDB/cluster"
	"MaybeDB/servers"
	"MaybeDB/utils"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

/*
 * 主节点重启后进行数据恢复
 */

//数据恢复到内存
func RecoveryInit() {

	//是否以集群方式部署
	isCluster := viper.GetInt("server.isCluster")
	//是否是主节点
	isMaster := viper.GetInt("server.isMaster")
	// 数据恢复策略
	recoveryStrategy := viper.GetInt("db.recoveryStrategy")

	//从本地持久化文件中获取数据
	if recoveryStrategy == 1 {
		//本地恢复数据
		fmt.Println("recovery from local")
		recoveryFromLocal()
	}

	//从集群其他健康的主节点获取数据
	if recoveryStrategy == 2 && isCluster == 1 && isMaster == 1 {
		//云端恢复数据
		fmt.Println("recovery from cluster")
		recoveryFromCluster()
	}
}

//从主节点集群中获取数据进行恢复工作
func recoveryFromCluster() {

	//获取一个健康的主节点实例，获取主节点上的Nacos元数据
	instance, err := cluster.NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
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

//从本地持久化文件中获取数据进行恢复工作
func recoveryFromLocal() {
	Read()
}
