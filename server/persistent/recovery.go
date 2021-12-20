package persistent

import (
	"MaybeDB/server/cluster"
	"MaybeDB/server/database"
	"MaybeDB/utils"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"strconv"
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

	//集群模式下，从节点会自动同步主节点数据，无需进行数据恢复操作
	if isCluster == 1 && isMaster == 0 {
		return
	}

	//从本地持久化文件中获取数据
	if recoveryStrategy == 1 {
		//本地恢复数据
		database.Loger.Println("Recovery from local")
		recoveryFromLocal()
	}

	//从集群其他健康的主节点获取数据
	if recoveryStrategy == 2 && isCluster == 1 && isMaster == 1 {
		//云端恢复数据
		database.Loger.Println("Recovery from cluster")
		recoveryFromCluster()
	}

	//recoveryStrategy为其他数值时不进行数据恢复操作
}

//从主节点集群中获取数据进行恢复工作
func recoveryFromCluster() {

	// SelectInstances 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, err := cluster.NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: "maybe-db-master",
		GroupName:   "MAYBE_DB_GROUP",             // 默认值DEFAULT_GROUP
		Clusters:    []string{"MAYBE_DB_CLUSTER"}, // 默认值DEFAULT
		HealthyOnly: true,
	})
	if err != nil {
		database.Loger.Println(err)
		return
	}

	//遍历实例列表，找到一个健康的实例，向其发出同步请求
	for _, instance := range instances {
		//如果是自己，跳过
		if instance.Ip == viper.GetString("server.ip") && instance.Port == viper.GetUint64("server.port") {
			continue
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
			continue
		}

		//解析数据到masterMap集合
		masterMap, err := utils.JsonToData(res)
		if err != nil {
			database.Loger.Println(err)
			continue
		}

		//将主节点map内的数据（masterMap）循环写入从节点map（DataMap）
		for key, value := range masterMap {
			database.DataMap.Store(key, value)
		}
		break
	}
}

//从本地持久化文件中获取数据进行恢复工作
func recoveryFromLocal() {
	Read()
}
