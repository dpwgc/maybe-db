# MaybeDB

## 基于Golang的简易key-value内存型数据库

`Golang` `Gin` `sync.Map`

***

### 实现功能

* 插入数据（支持string、int64、map、array类型的value）

* 根据key查找value

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

* clearServer `实时清理过期数据`

* connServer `客户端连接操作`

* dataServer `数据存储`

##### main.go 主函数

***

### 使用说明

* 填写application.yaml内的配置，运行main.go

***