# MaybeDB

## 基于Golang的简易key-value内存型数据库

`Golang` `Gin` `sync.Map` `DB`

***

### 实现功能

* 插入数据（仅支持string类型的value）

* 给数据设置过期时间

* 自动清除过期数据

* 删除数据

***

### 项目结构

##### config 配置类

* application.yaml `项目配置文件`

* config.go `项目配置文件加载`

##### models 实体类

* model.go `数据模板`

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