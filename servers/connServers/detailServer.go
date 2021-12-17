package connServers

import (
	"MaybeDB/servers"
	"github.com/gin-gonic/gin"
)

//根据key获取数据详情（展示数据的过期时间及数据类型）
func DetailGet(c *gin.Context) {

	key, _ := c.GetPostForm("key")

	value, _ := servers.DataMap.Load(key)

	c.JSON(0, gin.H{
		"code": 0,
		"data": value.(servers.Data),
	})
}

//获取数据详情列表（展示数据的过期时间及数据类型）
func DetailList(c *gin.Context) {

	resMap := make(map[string]servers.Data)

	var count int64 = 0

	servers.DataMap.Range(func(key, value interface{}) bool {
		resMap[key.(string)] = value.(servers.Data)
		count++
		return true
	})

	c.JSON(0, gin.H{
		"code":  0,
		"count": count,
		"data":  resMap,
	})
}

//获取当前数据总数
func Count(c *gin.Context) {

	var count int64 = 0

	servers.DataMap.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	c.JSON(0, gin.H{
		"code":  0,
		"count": count,
	})
}
