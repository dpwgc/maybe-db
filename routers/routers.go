package routers

import (
	"MaybeDB/clientConns"
	"MaybeDB/middlewares"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/**
 * 路由
 */

func SetupRouters() (r *gin.Engine) {

	r = gin.Default()
	r.Use(Cors())

	//客户端连接
	conn := r.Group("/Client")
	conn.Use(middlewares.SafeMiddleWare)
	{
		conn.POST("/Set", clientConns.Set)
		conn.GET("/Get", clientConns.Get)
		conn.DELETE("/Del", clientConns.Del)

		conn.GET("/List", clientConns.List)
		conn.GET("/ListByKeyword", clientConns.ListByKeyword)
		conn.GET("/ListByPrefix", clientConns.ListByPrefix)

		conn.GET("/DetailGet", clientConns.DetailGet)
		conn.GET("/DetailList", clientConns.DetailList)
		conn.GET("/DetailListByKeyword", clientConns.DetailListByKeyword)
		conn.GET("/DetailListByPrefix", clientConns.DetailListByPrefix)
		conn.GET("/Count", clientConns.Count)
	}
	return
}

//跨域设置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}
