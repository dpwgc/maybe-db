package cluster

import (
	"MaybeDB/servers"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var namingClient naming_client.INamingClient

func NacosInit() {

	clientConfig := constant.ClientConfig{
		NamespaceId:         "",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              viper.GetString("nacos.logDir"),   // 日志存储路径
		CacheDir:            viper.GetString("nacos.cacheDir"), // 缓存service信息的目录
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	//nacos注册中心地址配置
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      viper.GetString("nacos.ipAddr"), //nacos注册中心的ip
			ContextPath: viper.GetString("nacos.contextPath"),
			Port:        uint64(viper.GetInt("nacos.port")),
			Scheme:      viper.GetString("nacos.scheme"),
		},
	}

	// 将服务注册到nacos
	namingClient, _ = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	var matedata map[string]string
	var serviceName string
	//判断是否为主节点
	isMaster := viper.GetInt("server.isMaster")

	//如果该节点是主节点
	if isMaster == 1 {
		//将主节点ServiceName标记为maybe-db-master
		serviceName = "maybe-db-master"
		//将DataMap数据同步到nacos元数据
		matedata = map[string]string{"DataMap": servers.JsonCopyMap}
	}

	//如果该节点是从节点
	if isMaster == 0 {
		//将从节点ServiceName标记为maybe-db-slave
		serviceName = "maybe-db-slave"
		//从节点不对nacos元数据进行任何处理
		matedata = map[string]string{"Slave Node": "200"}
	}

	//向Nacos注册服务
	namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          viper.GetString("server.ip"),
		Port:        uint64(viper.GetInt("server.port")),
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    matedata,
		ClusterName: "MAYBE_DB",
		GroupName:   "MAYBE_DB_GROUP",
	})
}
