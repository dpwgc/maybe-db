package models

//数据模板

type Data struct {
	Content    interface{} //数据内容
	ExpireTime int64       //过期日期时间戳
}
