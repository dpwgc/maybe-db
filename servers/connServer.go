package servers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

/**
 * 客户端连接服务
 */

//插入数据
func Set(c *gin.Context) {

	key, _ := c.GetPostForm("key")
	value, _ := c.GetPostForm("value")
	expireTime, _ := c.GetPostForm("expireTime")

	expireTimeInt64, _ := strconv.ParseInt(expireTime, 10, 64)

	var expTime int64

	if expireTimeInt64 == 0 {
		expTime = 0
	} else {
		expTime = time.Now().Unix() + expireTimeInt64
	}

	data := Data{
		Content:    value,
		ExpireTime: expTime,
	}

	dataMap.Store(key, data)

	c.JSON(0, gin.H{
		"code": 0,
	})
}

//根据key获取数据
func Get(c *gin.Context) {

	key, _ := c.GetPostForm("key")

	value, _ := dataMap.Load(key)
	if value == nil {
		c.JSON(-1, gin.H{
			"code": -1,
		})
		return
	}

	c.JSON(0, gin.H{
		"code": 0,
		"data": value.(Data).Content,
	})
}

//根据key获取数据详情（展示数据的过期时间）
func DetailGet(c *gin.Context) {

	key, _ := c.GetPostForm("key")

	value, _ := dataMap.Load(key)

	c.JSON(0, gin.H{
		"code": 0,
		"data": value.(Data),
	})
}

//删除数据
func Del(c *gin.Context) {

	key, _ := c.GetPostForm("key")

	dataMap.Delete(key)

	c.JSON(0, gin.H{
		"code": 0,
	})
}

//获取数据列表
func List(c *gin.Context) {

	resMap := make(map[string]interface{})

	var count int64 = 0

	dataMap.Range(func(key, value interface{}) bool {
		resMap[key.(string)] = value.(Data).Content
		count++
		return true
	})

	c.JSON(0, gin.H{
		"code":  0,
		"count": count,
		"data":  resMap,
	})
}

//获取数据详情列表（展示数据的过期时间）
func DetailList(c *gin.Context) {

	resMap := make(map[string]Data)

	var count int64 = 0

	dataMap.Range(func(key, value interface{}) bool {
		resMap[key.(string)] = value.(Data)
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

	dataMap.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	c.JSON(0, gin.H{
		"code":  0,
		"count": count,
	})
}
