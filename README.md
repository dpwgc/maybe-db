# MaybeDB

## 基于Golang Gin整合Nacos的分布式key-value内存型数据库

`Golang` `Gin` `Nacos` `sync.Map`

***

### 集群实现功能

* 主从节点数据同步功能：将主节点数据定期更新Nacos的Matedata元数据空间，从节点从Nacos上获取主节点的元数据，将元数据中的map解析并覆盖本地数据。

### 单机实现功能

* 插入数据（value的类型：支持string、int64、map、int类型数组）

* 根据key查找数据

* 根据key中关键字获取数据列表

* 根据key的前缀获取数据列表

* 获取全部数据

* 统计数据总数

* 给数据设置过期时间

* 自动清除过期数据

* 删除数据

***

### 项目结构

##### clientConn 客户端连接操作

* setConn `存储数据相关操作`
* getConn `获取数据相关操作`
* delConn `删除数据相关`
* detailConn `获取数据详情相关操作`

##### cluster 集群相关

* dataSync `主从数据同步`
* nacos `Nacos注册中心连接`

##### config 配置类

* application.yaml `项目配置文件`

* config.go `项目配置文件加载`

##### middlewares 中间件

* safeMiddleware `访问密钥验证`

##### routers 路由

* routers.go `路由配置`

##### servers 服务层

* clearServer `实时清理过期数据`

* dataServer `数据存储`

##### main.go 主函数

***

### 使用说明

* 填写application.yaml内的配置。

* 运行项目：

```
直接运行main.go(调试)
打包成exe运行(windows)
打包成二进制文件运行(linux)
```

***

### 部署方式

* 在一台服务器上部署Nacos，填写各节点application.yaml配置文件中的Nacos配置信息。

* 采用一主多从方式部署项目,主节点写数据，从节点读数据。

* application.yaml主从节点配置

```
# 主节点配置
isMaster: 1

# 从节点配置
isMaster: 0
```

* 将主从节点打包好，分别部署上服务器即可。
