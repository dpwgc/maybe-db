package diskStorage

import (
	"MaybeDB/servers"
	"MaybeDB/utils"
	"encoding/csv"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

/*
 * 持久化文件读写
 */

var wFile *os.File
var rFile *os.File
var path string

//初始化文件
func FileInit() {

	path = viper.GetString("db.persistentPath")
	//判断持久化文件是否存在
	_, err := os.Stat(path)
	if err != nil {
		//创建持久化文件
		fmt.Println("create" + path)
		_, _ = os.Create(path)
	}
}

//写入持久化文件
func Write() {
	var err error
	//写文件，设置为只写、覆盖，权限设置为777
	wFile, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)

	writer := csv.NewWriter(wFile)

	//将数据以json字符串形式存入持久化文件
	jsonStr := string(servers.PersCopyByte)
	err = writer.Write([]string{jsonStr})
	if err != nil {
		fmt.Print("Write:")
		fmt.Println(err)
		return
	}
	writer.Flush()
	//关闭文件流
	err = wFile.Close()
	if err != nil {
		fmt.Print("Write Close:")
		fmt.Println(err)
	}
}

//加载持久化文件内的数据到内存中
func Read() {
	var err error
	//读文件，设置为只读，权限设置为777
	rFile, err = os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		fmt.Print("Read OpenFile:")
		fmt.Println(err)
		return
	}
	reader := csv.NewReader(rFile)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		fmt.Print("Read ReadAll:")
		fmt.Println(err)
		return
	}
	if len(record) == 0 {
		fmt.Println("EOF")
		return
	}
	//解析本地持久化文件的数据到localMap
	localMap, err := utils.JsonToData(record[0][0])
	if err != nil {
		fmt.Print("Read JsonToData:")
		fmt.Println(err)
		return
	}

	//将本地持久化文件数据（localMap）循环写入从节点map（DataMap）
	for key, value := range localMap {
		servers.DataMap.Store(key, value)
	}
	err = rFile.Close()
	if err != nil {
		fmt.Print("Read Close:")
		fmt.Println(err)
	}
}
