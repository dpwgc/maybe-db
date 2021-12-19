package main

import (
	"MaybeDB/cluster"
	"MaybeDB/config"
	"MaybeDB/diskStorage"
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

	//加载配置文件
	config.ConfigInit()

	//加载日志记录模块
	servers.LogInit()

	//加载文件操作模块
	diskStorage.FileInit()

	//是否开启持久化
	isPersistent := viper.GetInt("db.isPersistent")
	if isPersistent == 1 {
		diskStorage.PersInit()
	}

	//是否以集群方式部署
	isCluster := viper.GetInt("server.isCluster")
	if isCluster == 1 {
		cluster.NacosInit()
		cluster.SyncInit()
	}

	//加载数据恢复模块
	diskStorage.RecoveryInit()

	//初始化清理模块
	servers.ClearInit()

	//设置路由
	r := routers.SetupRouters()

	//获取端口号
	port := viper.GetString("server.port")
	_ = r.Run(":" + port)
}
