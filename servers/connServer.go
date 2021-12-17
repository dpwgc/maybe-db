package servers

import (
	"MaybeDB/utils"
	"encoding/json"
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
	valueType, _ := c.GetPostForm("valueType")
	expireTime, _ := c.GetPostForm("expireTime")

	valueTypeInt, err := strconv.Atoi(valueType)
	if err != nil {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "type conversion error",
		})
		return
	}

	expireTimeInt64, err := strconv.ParseInt(expireTime, 10, 64)
	if err != nil {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "type conversion error",
		})
		return
	}

	var expTime int64

	if expireTimeInt64 == 0 {
		expTime = 0
	} else {
		expTime = time.Now().Unix() + expireTimeInt64
	}

	var content interface{}

	switch valueTypeInt {

	//string类型的值
	case 1:
		content = value
		break

	//int64类型的值
	case 2:
		//string转int
		content, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			c.JSON(-1, gin.H{
				"code": -1,
				"msg":  "type conversion error",
			})
			return
		}
		break

	//map类型的值（value为json字符串，例：`{"id": 1, "text": "hello"}`）
	case 3:
		//json转map
		content, err = utils.JsonToMap(value)
		if err != nil {
			c.JSON(-1, gin.H{
				"code": -1,
				"msg":  "type conversion error",
			})
			return
		}
		break

	//array类型的值（value为json字符串，例：`[1, 2, 3, 4, 5]`）
	case 4:
		//将json字符串解析为数组
		err = json.Unmarshal([]byte(value), &content)
		if err != nil {
			c.JSON(-1, gin.H{
				"code": -1,
				"msg":  "type conversion error",
			})
			return
		}
		break

	//无效类型
	default:
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "there is no such type",
		})
		return
	}

	//生成数据模板
	data := Data{
		Content:    content,
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
			"msg":  "value cannot be empty",
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
