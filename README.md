# MaybeDB

## 基于Golang整合Nacos的分布式键值型内存数据库

`Golang` `Gin` `Nacos` `sync.Map`

***

### 集群实现功能

* 主从节点读写分离：主节点负责插入/删除数据，从节点负责查询数据。

* 高可用：支持一主多从集群部署或多主多从集群部署，单个主节点宕机重启后，可拉取其他健康的主节点的数据或从本地持久化文件中提取数据来进行数据恢复。

* 主节点之间数据同步：在任意一个主节点进行数据插入/删除操作，该操作都会扩散同步到其他主节点。

* 主从节点数据同步：从节点从Nacos上获取任意一个健康的主节点实例，像该主节点实例发出同步请求，并将请求结果解析同步到本地。

### 单机实现功能

* 数据持久化与重启数据恢复

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

* sync `主从节点数据同步`

* recovery `主节点数据恢复`

* nacos `Nacos注册中心连接`

##### config 配置类

* application.yaml `项目配置文件`

* config.go `项目配置文件加载`

##### diskStorage 持久化到硬盘

* fileRW `文件读写操作`

* persistent `数据持久化`

* recovery `数据恢复`

##### middlewares 中间件

* safeMiddleware `访问密钥验证`

##### routers 路由

* routers.go `路由配置`

##### servers 服务层

* clearServer `实时清理过期数据`

* dataServer `数据存储`

##### utils 工具类

* httpUtil.go `http请求工具`

* jsonUtil `Json字符串转换工具`

##### main.go 主函数

***

### 打包方式

* 填写application.yaml内的配置。

* 运行项目：

```
GoLand直接运行main.go(调试)

打包成exe运行(windows部署)
GoLand终端cd到项目根目录，执行go build main.go命令，生成main.exe文件

打包成二进制文件运行(linux部署)
cmd终端cd到项目根目录，依次执行下列命令：
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go
生成main文件
```

***

### 部署说明

##### 集群部署

* 在一台服务器上部署Nacos，添加命名空间(命名空间Id及命名空间名称设置为：maybe-db),填写各节点application.yaml配置文件中的Nacos配置信息。

* 采用一主多从/多主多从方式部署项目,主节点写数据，从节点读数据。

* application.yaml主从节点配置

```
# 主节点配置
isCluster: 1
isMaster: 1

# 从节点配置
isCluster: 1
isMaster: 0
```

* 将主从节点打包好，分别上传至服务器运行即可。

```
Windows
/maybe-db                 # 应用文件根目录
    main.exe              # 打包后的exe文件
    /config               # 配置目录
        application.yaml  # 配置文件
    /cache                # Nacos缓存目录
    /log                  # Nacos日志目录
    DataMap.csv           # 持久化文件
    
Linux
/maybe-db                 # 应用文件根目录
    main                  # 打包后的二进制文件(程序后台执行:setsid ./main)
    /config               # 配置目录
        application.yaml  # 配置文件
    /cache                # Nacos缓存目录
    /log                  # Nacos日志目录
    DataMap.csv           # 持久化文件
```

##### 单机部署

* 无需部署Nacos，将application.yaml配置文件中的isCluster设置为0，直接打包部署到服务器。

```
# 是否以集群方式部署（1:是，0:否）
isCluster: 0
# 是否为主节点（1:是，0:否）
isMaster: 1
```

##### 数据持久化

* 集群模式下，从节点会自动同步主节点数据，无需进行数据持久化/数据恢复

```
# 是否开启持久化（1:是，0:否）
isPersistent: 1
# 数据持久化操作的时间间隔（单位：秒）
persistentTime: 5
# 持久化文件保存路径
persistentPath: ./DataMap.csv
```