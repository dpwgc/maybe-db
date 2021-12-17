package connServers

import (
	"MaybeDB/servers"
	"github.com/gin-gonic/gin"
	"strings"
)

//根据key获取数据
func Get(c *gin.Context) {

	key, _ := c.GetPostForm("key")

	value, _ := servers.DataMap.Load(key)
	if value == nil {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "value cannot be empty",
		})
		return
	}

	c.JSON(0, gin.H{
		"code": 0,
		"data": value.(servers.Data).Content,
	})
}

//获取全部数据列表
func List(c *gin.Context) {

	resMap := make(map[string]interface{})

	var count int64 = 0

	servers.DataMap.Range(func(key, value interface{}) bool {
		resMap[key.(string)] = value.(servers.Data).Content
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
func ListByKeyword(c *gin.Context) {

	keyword, _ := c.GetPostForm("keyword")

	resMap := make(map[string]interface{})

	var count int64 = 0

	servers.DataMap.Range(func(key, value interface{}) bool {
		//如果查找到关键字
		if len(strings.Split(key.(string), keyword)) > 1 {
			resMap[key.(string)] = value.(servers.Data).Content
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
func ListByPrefix(c *gin.Context) {

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
		resMap[key.(string)] = value.(servers.Data).Content
		count++
		return true
	})

	c.JSON(0, gin.H{
		"code":  0,
		"count": count,
		"data":  resMap,
	})
}
