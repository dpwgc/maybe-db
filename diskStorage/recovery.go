package diskStorage

import (
	"MaybeDB/cluster"
	"MaybeDB/servers"
	"MaybeDB/utils"
	"encoding/json"
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
		fmt.Println("diskStorage from local")
		recoveryFromLocal()
	}

	//从集群其他健康的主节点获取数据
	if recoveryStrategy == 2 && isCluster == 1 && isMaster == 1 {
		//云端恢复数据
		fmt.Println("diskStorage from cluster")
		recoveryFromCluster()
	}

	//recoveryStrategy为其他数值时不进行数据恢复操作
}

//从主节点集群中获取数据进行恢复工作
func recoveryFromCluster() {

	//获取一个健康的主节点实例
	instance, err := cluster.NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "maybe-db-master",
		GroupName:   "MAYBE_DB_GROUP",
		Clusters:    []string{"MAYBE_DB_CLUSTER"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//设置请求头
	header := make(map[string]string, 1)
	//访问密钥
	header["secretKey"] = viper.GetString("db.secretKey")

	//向该主节点请求数据
	url := fmt.Sprintf("%s%s%s%s%s", "http://", instance.Ip, ":", strconv.Itoa(int(instance.Port)), "/Sync/GetMasterData")
	resStr, err := utils.Get(url, header)

	var res string
	err = json.Unmarshal([]byte(resStr), &res)
	if err != nil {
		fmt.Println(err)
	}

	//解析数据到masterMap集合
	masterMap, err := utils.JsonToData(res)
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
