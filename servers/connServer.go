package servers

import (
	"MaybeDB/models"
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

	data := models.Data{
		Content:    value,
		ExpireTime: time.Now().Unix() + expireTimeInt64,
	}

	dataMap.Store(key, data)

	c.JSON(0, gin.H{
		"code": 0,
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
		resMap[key.(string)] = value.(models.Data).Content
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
