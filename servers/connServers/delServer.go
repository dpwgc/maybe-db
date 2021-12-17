package connServers

import (
	"MaybeDB/servers"
	"github.com/gin-gonic/gin"
)

//删除数据
func Del(c *gin.Context) {

	key, _ := c.GetPostForm("key")

	servers.DataMap.Delete(key)

	c.JSON(0, gin.H{
		"code": 0,
	})
}
