package utils

import (
	"MaybeDB/servers"
	"encoding/json"
)

//json字符串转map
func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

//json字符串转Data结构体
func JsonToData(jsonStr string) (map[string]servers.Data, error) {
	m := make(map[string]servers.Data)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
