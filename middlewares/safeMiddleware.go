package middlewares

import (
	"MaybeDB/servers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//安全防护中间件-访问密钥验证
func SafeMiddleWare(c *gin.Context) {

	secretKey := c.GetHeader("secretKey")

	if secretKey != viper.GetString("db.secretKey") {

		sourceIp := c.ClientIP()
		accessUri := c.Request.URL.Path

		servers.Loger.Println("Illegal access , sourceIp: " + sourceIp + " , accessUri: " + accessUri)

		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "key matching error",
		})
		c.Abort()
	}
	return
}
