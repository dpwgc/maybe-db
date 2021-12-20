package main

import (
	"MaybeDB/config"
	"MaybeDB/router"
	"MaybeDB/server/cluster"
	"MaybeDB/server/database"
	"MaybeDB/server/persistent"
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
	database.LogInit()

	//加载文件操作模块
	persistent.FileInit()

	//加载数据恢复模块
	persistent.RecoveryInit()

	//加载持久化模块
	persistent.PersInit()

	//加载Nacos注册中心连接模块
	cluster.NacosInit()

	//加载主从数据同步模块
	cluster.SyncInit()

	//初始化清理模块
	database.ClearInit()

	//设置路由
	r := router.SetupRouters()

	//获取端口号
	port := viper.GetString("server.port")
	_ = r.Run(":" + port)
}
