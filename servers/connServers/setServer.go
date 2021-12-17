package connServers

import (
	"MaybeDB/servers"
	"MaybeDB/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

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
			"msg":  err,
		})
		return
	}

	expireTimeInt64, err := strconv.ParseInt(expireTime, 10, 64)
	if err != nil {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  err,
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
				"msg":  err,
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
				"msg":  err,
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
				"msg":  err,
			})
			return
		}
		break

	//无效类型
	default:
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  err,
		})
		return
	}

	//生成数据模板
	data := servers.Data{
		Content:    content,
		ExpireTime: expTime,
	}

	servers.DataMap.Store(key, data)

	c.JSON(0, gin.H{
		"code": 0,
	})
}
