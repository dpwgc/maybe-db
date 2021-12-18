package utils

import (
	"MaybeDB/servers"
	"encoding/csv"
	"fmt"
	"os"
)

var csvFile *os.File
var path = "./DataMap.csv"

func init() {
	//判断持久化文件是否存在
	_, err := os.Stat(path)
	if err != nil {
		//创建持久化文件
		_, _ = os.Create(path)
	}
	//打开文件，设置为覆盖写入，权限设置为777
	csvFile, err = os.OpenFile(path, os.O_TRUNC, 0777)
}

//覆盖写入持久化文件
func Write() {
	writer := csv.NewWriter(csvFile)
	err := writer.Write([]string{servers.PersCopyJson})
	if err != nil {
		panic(err)
	}
	writer.Flush()
}

//加载持久化文件内的数据到内存中
func Read() {
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	record, err := reader.Read()
	if err != nil {
		panic(err)
	}
	//解析本地持久化文件数据到localMap
	localMap, err := JsonToData(record[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	//将本地持久化文件数据（localMap）循环写入从节点map（DataMap）
	for key, value := range localMap {
		servers.DataMap.Store(key, value)
	}
}
