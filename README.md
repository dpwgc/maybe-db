# MaybeDB

## 基于Golang的简易key-value内存型数据库

`Golang` `Gin` `sync.Map`

***

### 实现功能

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

##### config 配置类

* application.yaml `项目配置文件`

* config.go `项目配置文件加载`

##### middlewares 中间件

* 访问密钥验证

##### routers 路由

* routers.go `路由配置`

##### servers 服务层

* clientConnServers `客户端连接相关操作`

```
setServer 存储数据相关操作
getServer 获取数据相关操作
delServer 删除数据相关
detailServer 获取数据详情相关操作
```

* clearServer `实时清理过期数据`

* dataServer `数据存储`

##### main.go 主函数

***

### 使用说明

* 填写application.yaml内的配置，运行main.go
