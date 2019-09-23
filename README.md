# pool
[![GoDoc](http://godoc.org/github.com/silenceper/pool?status.svg)](http://godoc.org/github.com/silenceper/pool)

Golang 实现的连接池


## 功能：

- 连接池中连接类型为`interface{}`，使得更加通用
- 连接的最大空闲时间，超时的连接将关闭丢弃，可避免空闲时连接自动失效问题
- 支持用户设定 ping 方法，检查连接的连通性，无效的连接将丢弃
- 使用channel处理池中的连接，高效
- 最大连接数限定 避免可以无限创建连接

## 基本用法

```go

//factory 创建连接的方法
factory := func() (interface{}, error) { return net.Dial("tcp", "127.0.0.1:8000") }

//close 关闭连接的方法
close := func(v interface{}) error { return v.(net.Conn).Close() }

//ping 检测连接的方法
//ping := func(v interface{}) error { return nil }

//创建一个连接池： 初始化5，最大连接30
poolConfig := &pool.PoolConfig{
	InitialCap: 5,
	MaxCap:     30,
	Factory:    factory,
	Close:      close,
	//Ping:       ping,
	//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
	IdleTimeout: 15 * time.Second,
}
p, err := pool.NewChannelPool(poolConfig)
if err != nil {
	fmt.Println("err=", err)
}

//从连接池中取得一个连接
v, err := p.Get()

//do something
//conn=v.(net.Conn)

//将连接放回连接池中 取出的连接 要么Put回去或者Close掉 否则限制数将不准确
p.Put(v)

//查看当前空闲连接的数量
size := p.Len()

//查看当前还可创建连接的数量
remain := p.Remain()

//使用中的连接数量
//Max-p.Len()-p.Remain()

//释放连接池中的所有连接
p.Release()

```

## License

The MIT License (MIT) - see LICENSE for more details
