package servers

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

/**
 * 日志记录
 */

var Loger *log.Logger

func LogInit() {

	ip := viper.GetString("server.ip")
	port := viper.GetString("server.port")

	file := "./log/maybe-db-" + time.Now().Format("2006-01") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		//创建log目录
		err = os.Mkdir("./log", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	Loger = log.New(logFile, "["+ip+":"+port+"]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
}
