package servers

import (
	"sync"
)

/**
 * 数据存储
 */

func init() {
	CopyMap = make(map[string]interface{})
}

//使用sync.Map存放数据
var DataMap sync.Map

//主从复制时使用（cluster/dataSync.go）
var CopyMap map[string]interface{}
var ByteCopyMap []byte
var JsonCopyMap string

//数据模板
type Data struct {
	Content     interface{} //数据内容
	ContentType int         //数据类型（1:string，2:int64，3:map，4:array）
	ExpireTime  int64       //过期日期时间戳（为0时表示永不删除）
}
