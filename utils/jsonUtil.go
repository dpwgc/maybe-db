package utils

import (
	"MaybeDB/server/database"
	"encoding/json"
)

/*
 * JSON格式转换工具
 */

//json字符串转map
func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		database.Loger.Println(err)
		return nil, err
	}
	return m, nil
}

//json字符串转Data结构体
func JsonToData(jsonStr string) (map[string]database.Data, error) {
	m := make(map[string]database.Data)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		database.Loger.Println(err)
		return nil, err
	}
	return m, nil
}

//map转json字符串
func MapToJson(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}
