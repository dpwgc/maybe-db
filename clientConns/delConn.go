package clientConns

import (
	"MaybeDB/servers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//删除数据
func Del(c *gin.Context) {

	isCluster := viper.GetInt("server.isCluster")
	isMaster := viper.GetInt("server.isMaster")

	//集群模式下，只有主节点可以写入/删除数据，从节点只负责读取数据
	if isMaster == 0 && isCluster == 1 {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "only the master node can delete data",
		})
		return
	}

	key, _ := c.GetPostForm("key")

	servers.DataMap.Delete(key)

	c.JSON(0, gin.H{
		"code": 0,
	})
}
