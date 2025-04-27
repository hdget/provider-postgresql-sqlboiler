### MySQL配置
mysql配置必须在`[sdk]`段落下, 同时必须包含`[sdk.log]`, 
日志能力`[sdk.log]`的配置是最基本需求的配置，不管使用sdk的时候是否具有其他能力, 日志配置信息必须要包含


1. 数据库连接类型

- default: 缺省的数据库
- master:  主数据库
- slave:   从数据库
- other:   额外数据库

2. 配置示例
```
[sdk.postgresql]
    [sdk.postgresql.default]
        database = "postgres"      <--- 要连接的数据库的名字
        user = "postgres"           <--- 数据库连接用户名
        password = "password"   <--- 数据库连接密码
        host = "127.0.0.1"      <--- 数据库主机
        port = 5432             <--- 数据库端口
    [sdk.postgresql.master]
        database = "postgres"
        user = "postgres"
        password = "password"
        host = "127.0.0.1"
        port = 5432
    [[sdk.mysql.slaves]]
        database = "postgres"
        user = "postgres"
        password = "password"
        host = "127.0.0.1"
        port = 5432
    [[sdk.mysql.slaves]]
        ...
    [[sdk.mysql.items]]
        name = "extra1"       <--- 额外的数据库连接的名字，用来区分不同的数据库连接
        database = "postgres"
        user = "postgres"
        password = "password"
        host = "127.0.0.1"
        port = 5432
    [[sdk.mysql.items]]
        name = "extra2"
        database = "postgres"
        user = "postgres"
        password = "password"
        host = "127.0.0.1"
        port = 5432
```

> 1. 数据库能力配置中的`name`只有在使用items的时候才需要指定
> 2. 在配置slave和other类型的时候需要注意用`[[`和`]]`, 另外名字后需要加`s`, e,g: slave的配置为`[[sdk.db.slaves]]`
 
### PostgreSQL使用指南

##### 获取数据库连接
- 获取缺省数据库连接: `sdk.Db.My()`
- 获取主数据库连接: `sdk.Db.Master()`
- 获取第几个从数据库连接: `sdk.Db.Slave(int)`
- 获取指定名字的数据库连接: `sdk.Db.By(string)`

#### 数据库接口

- 将占位符的查询语句转换成数据库driver对应的bindvar类型

    `func Rebind(query string) string`

- 查询获取一行并且将结果集中的指定字段的值防止到结果中，参考：https://pkg.go.dev/github.com/jmoiron/sqlx#Get

  `func Get(q Queryer, dest interface{}, query string, args ...interface{}) error`
  
- 查询数据库并返回*sqlx.Rows, query字符串中的`?`占位符会被提供的参数填入
  
  `func Queryx(query string, args ...interface{}) (*Rows, error)`
  
- 执行query语句

    `func Exec(query string, args ...interface{}) sql.Result`
