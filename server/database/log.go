package database

import (
	"log"
	"os"
	"time"
)

/**
 * 日志记录
 */

var Loger *log.Logger

func LogInit() {

	file := "./log/maybe-db-" + time.Now().Format("2006-01") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		//创建log目录
		err = os.Mkdir("./log", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	Loger = log.New(logFile, "", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
}
