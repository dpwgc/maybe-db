package main

import (
	"MaybeDB/config"
	"MaybeDB/registry"
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
	config.InitConfig()

	registry.NacosInit()
	registry.SyncInit()

	//初始化清理模块
	servers.InitClear()

	//设置路由
	r := routers.SetupRouters()

	//获取端口号
	port := viper.GetString("server.port")
	_ = r.Run(":" + port)
}
