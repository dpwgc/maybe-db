package main

import (
	"MaybeDB/cluster"
	"MaybeDB/config"
	"MaybeDB/recovery"
	"MaybeDB/routers"
	"MaybeDB/servers"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "net/http"
)

/**
 * main
 */

func main() {

	//加载配置
	config.ConfigInit()

	//数据恢复
	recovery.RecoveryInit()

	//是否开启持久化
	isPersistent := viper.GetInt("db.isPersistent")
	if isPersistent == 1 {
		recovery.PersInit()
	}

	//是否以集群方式部署
	isCluster := viper.GetInt("server.isCluster")
	if isCluster == 1 {
		cluster.NacosInit()
		cluster.SyncInit()
	}

	//初始化清理模块
	servers.ClearInit()

	//设置路由
	r := routers.SetupRouters()

	//获取端口号
	port := viper.GetString("server.port")
	_ = r.Run(":" + port)
}
