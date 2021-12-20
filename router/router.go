package router

import (
	"MaybeDB/clientConn"
	"MaybeDB/middleware"
	"MaybeDB/server/cluster"
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
	client := r.Group("/Client")
	client.Use(middleware.SafeMiddleWare)
	{
		client.POST("/Set", clientConn.Set)
		client.GET("/Get", clientConn.Get)
		client.DELETE("/Del", clientConn.Del)

		client.GET("/List", clientConn.List)
		client.GET("/ListByKeyword", clientConn.ListByKeyword)
		client.GET("/ListByPrefix", clientConn.ListByPrefix)

		client.GET("/DetailGet", clientConn.DetailGet)
		client.GET("/DetailList", clientConn.DetailList)
		client.GET("/DetailListByKeyword", clientConn.DetailListByKeyword)
		client.GET("/DetailListByPrefix", clientConn.DetailListByPrefix)
		client.GET("/Count", clientConn.Count)
	}
	//主从复制连接
	sync := r.Group("/Sync")
	sync.Use(middleware.SafeMiddleWare)
	{
		sync.GET("/GetMasterData", cluster.GetMasterData)
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
