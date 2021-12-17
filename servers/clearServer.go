package servers

import (
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

	DataMap.Range(func(key, value interface{}) bool {
		// 如果数据到期
		if value.(Data).ExpireTime <= nowTime && value.(Data).ExpireTime != 0 {
			//删除该数据
			DataMap.Delete(key)
		}
		return true
	})
}
