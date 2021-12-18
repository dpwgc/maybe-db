package servers

import (
	"sync"
)

/**
 * 数据存储
 */

func init() {
	SyncCopyMap = make(map[string]interface{})
	PersCopyMap = make(map[string]interface{})
}

//使用sync.Map存放数据
var DataMap sync.Map

//主从复制时使用（cluster/sync.go）
var SyncCopyMap map[string]interface{}
var SyncCopyByte []byte
var SyncCopyJson string

//持久化时使用（servers/persistent.go）
var PersCopyMap map[string]interface{}
var PersCopyByte []byte
var PersCopyJson string

//数据模板
type Data struct {
	Content     interface{} //数据内容
	ContentType int         //数据类型（1:string，2:int64，3:map，4:array）
	ExpireTime  int64       //过期日期时间戳（为0时表示永不删除）
}
