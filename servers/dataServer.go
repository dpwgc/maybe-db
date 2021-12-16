package servers

import (
	"sync"
)

/**
 * 数据存储
 */

//使用sync.Map存放数据
var dataMap sync.Map

//数据模板
type Data struct {
	Content    interface{} //数据内容
	ExpireTime int64       //过期日期时间戳（为0时表示永不删除）
}
