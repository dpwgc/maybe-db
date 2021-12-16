package servers

import (
	"MaybeDB/models"
	"time"
)

/**
 * 过期数据清理
 */

func InitClear() {
	go func() {
		for {
			time.Sleep(time.Second * 1)
			clear()
		}
	}()
}

func clear() {

	//获取当前时间戳
	nowTime := time.Now().Unix()

	dataMap.Range(func(key, value interface{}) bool {
		// 如果数据到期
		if value.(models.Data).ExpireTime <= nowTime {
			//删除该数据
			dataMap.Delete(key)
		}
		return true
	})
}
