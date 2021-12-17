package registry

import (
	"MaybeDB/servers"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"strings"
)

var namingClient naming_client.INamingClient

func NacosInit() {

	clientConfig := constant.ClientConfig{
		NamespaceId:         "",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "log",
		CacheDir:            "cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

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

	namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          viper.GetString("server.ip"),
		Port:        uint64(viper.GetInt("server.port")),
		ServiceName: "maybe-db",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"DataMap": servers.JsonCopyMap},
		ClusterName: "MAYBE_DB",
		GroupName:   "MAYBE_DB_GROUP",
	})

	//获取nacos存在服务的信息
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "maybe-db",
		GroupName:   "MAYBE_DB_GROUP",
		Clusters:    []string{"MAYBE_DB"},
	})
	fmt.Println(instance)
	if err != nil {
		fmt.Println(err)
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	//获取配置
	configClientcontent, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "maybe-db",
		Group:  "MAYBE_DB_GROUP"})
	fmt.Println(configClientcontent)
	if err != nil {
		fmt.Println(err)
	}

	//解析yml格式的数据完成初始化
	v := viper.New()
	v.SetConfigType("yaml")
	v.ReadConfig(strings.NewReader(configClientcontent))
	logoText := v.GetString("logo")
	fmt.Println(logoText)
}
