package cluster

import (
	"MaybeDB/servers"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var NamingClient naming_client.INamingClient

func NacosInit() {

	clientConfig := constant.ClientConfig{
		NamespaceId:         "maybe-db",
		TimeoutMs:           viper.GetUint64("nacos.TimeoutMs"),         //连接超时时间
		NotLoadCacheAtStart: viper.GetBool("nacos.notLoadCacheAtStart"), //账户
		LogDir:              viper.GetString("nacos.logDir"),            //日志存储路径
		CacheDir:            viper.GetString("nacos.cacheDir"),          //缓存service信息的目录
		RotateTime:          viper.GetString("nacos.rotateTime"),        //日志轮转周期
		MaxAge:              viper.GetInt64("nacos.maxAge"),             //日志最大文件数
		LogLevel:            viper.GetString("nacos.logLevel"),          //日志级别
		Username:            viper.GetString("nacos.username"),          //账户
		Password:            viper.GetString("nacos.password"),          //密码
	}

	//nacos注册中心地址配置
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      viper.GetString("nacos.ipAddr"), //nacos注册中心的ip
			ContextPath: viper.GetString("nacos.contextPath"),
			Port:        viper.GetUint64("nacos.port"),
			Scheme:      viper.GetString("nacos.scheme"),
		},
	}

	// 将服务注册到nacos
	NamingClient, _ = clients.NewNamingClient(
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
		matedata = map[string]string{"DataMap": servers.SyncCopyJson}
	}

	//如果该节点是从节点
	if isMaster == 0 {
		//将从节点ServiceName标记为maybe-db-slave
		serviceName = "maybe-db-slave"
		//从节点不对nacos元数据进行任何处理
		matedata = map[string]string{"Slave Node": "200"}
	}

	//向Nacos注册服务
	NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          viper.GetString("server.ip"),
		Port:        uint64(viper.GetInt("server.port")),
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    matedata,
		ClusterName: "MAYBE_DB_CLUSTER",
		GroupName:   "MAYBE_DB_GROUP",
	})
}
