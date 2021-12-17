package connServers

import (
	"MaybeDB/servers"
	"github.com/gin-gonic/gin"
	"strings"
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

//获取全部数据详情列表（展示数据的过期时间及数据类型）
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

//根据关键字获取数据列表
func DetailListByKeyword(c *gin.Context) {

	keyword, _ := c.GetPostForm("keyword")

	resMap := make(map[string]interface{})

	var count int64 = 0

	servers.DataMap.Range(func(key, value interface{}) bool {
		//如果查找到关键字
		if len(strings.Split(key.(string), keyword)) > 1 {
			resMap[key.(string)] = value.(servers.Data)
			count++
		}
		return true
	})

	c.JSON(0, gin.H{
		"code":  0,
		"count": count,
		"data":  resMap,
	})
}

//根据key前缀获取数据列表
func DetailListByPrefix(c *gin.Context) {

	prefix, _ := c.GetPostForm("prefix")
	size := len(prefix)

	resMap := make(map[string]interface{})

	var count int64 = 0

	servers.DataMap.Range(func(key, value interface{}) bool {

		k := key.(string)
		//前缀匹配
		for i := 0; i < size; i++ {
			//发现不匹配字符
			if k[i] != prefix[i] {
				//跳过该数据
				return true
			}
		}
		//匹配成功
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
