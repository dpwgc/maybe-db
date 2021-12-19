package cluster

import (
	"MaybeDB/servers"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
 * 主从同步，主节点操作
 */

//该主节点的DataMap数据获取接口，用于提供给从节点进行数据同步操作
func GetMasterData(c *gin.Context) {

	//以Json字符串形式返回主节点的全部数据
	res := string(servers.SyncCopyByte)

	c.String(http.StatusOK, fmt.Sprintln(res))
}

//复制该主节点的本地数据
func copyDataMap() {

	copyMap := make(map[string]interface{})

	servers.DataMap.Range(func(key, value interface{}) bool {
		copyMap[key.(string)] = value
		return true
	})
	//将SyncCopyMap转为字节数组类型SyncCopyByte
	servers.SyncCopyByte, _ = json.Marshal(copyMap)
}
