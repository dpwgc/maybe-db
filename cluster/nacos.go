package cluster

import (
	"MaybeDB/servers"
	"MaybeDB/utils"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var NamingClient naming_client.INamingClient

func NacosInit() {

	//Nacos详细配置说明：https://github.com/nacos-group/nacos-sdk-go/blob/master/README_CN.md
	clientConfig := constant.ClientConfig{
		NamespaceId: viper.GetString("nacos.namespaceId"), //命名空间Id
		TimeoutMs:   viper.GetUint64("nacos.TimeoutMs"),   //连接超时时间

		Endpoint:  viper.GetString("nacos.endpoint"),
		RegionId:  viper.GetString("nacos.regionId"),
		AccessKey: viper.GetString("nacos.accessKey"),
		SecretKey: viper.GetString("nacos.secretKey"),
		OpenKMS:   viper.GetBool("nacos.openKMS"),

		UpdateThreadNum:      viper.GetInt("nacos.updateThreadNum"),
		NotLoadCacheAtStart:  viper.GetBool("nacos.notLoadCacheAtStart"),
		UpdateCacheWhenEmpty: viper.GetBool("nacos.updateCacheWhenEmpty"),

		LogDir:     "log",                               //日志存储路径
		CacheDir:   "cache",                             //缓存service信息的目录
		RotateTime: viper.GetString("nacos.rotateTime"), //日志轮转周期
		MaxAge:     viper.GetInt64("nacos.maxAge"),      //日志最大文件数
		LogLevel:   viper.GetString("nacos.logLevel"),   //日志级别

		Username: viper.GetString("nacos.username"), //账户
		Password: viper.GetString("nacos.password"), //密码
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

	var serviceName string
	//判断是否为主节点
	isMaster := viper.GetInt("server.isMaster")

	//如果该节点是主节点
	if isMaster == 1 {
		//将主节点ServiceName标记为maybe-db-master
		serviceName = "maybe-db-master"
	}

	//如果该节点是从节点
	if isMaster == 0 {
		//将从节点ServiceName标记为maybe-db-slave
		serviceName = "maybe-db-slave"
	}

	//将配置信息上传到元数据空间
	matedata := map[string]string{"Config": utils.MapToJson(viper.AllSettings())}

	//向Nacos注册服务
	_, err := NamingClient.RegisterInstance(vo.RegisterInstanceParam{
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
	if err != nil {
		servers.Loger.Println(err)
		return
	}
}
