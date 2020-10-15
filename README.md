## Vex

> 基于 mit 协议开源

一个通用网络通信框架，可以作为参考案例进行交流学习哈哈。因为请求和响应发送的时候是分开的，所以不支持并发请求，这个要改也很简单，只要一次性发送请求和相应即可。

不过这个要看应用场景，如果发送的数据太大，容易造成内存的浪费和分配，所以要看实际应用。

ANBF 协议：

```anbf
HEADER = VERSION CMDLENGTH CMD ARGSLENGTH ; 头部，主要是协议版本号，指令长度和指令，参数长度
BODY = *{ARGLENGTH ARG} ; 正文，主要是参数列表，每个参数都以长度开头
VERSION = OCTET ; 协议版本号，OCTET 是指 0x00-0xFF
CMDLENGTH = 2OCTET ; 指令长度，0x00 - 0xFFFF
CMD = *OCTET ; 指令，字节数组形式
ARGSLENGTH = 4OCTET ; 参数个数，4 个字节
ARGLENGTH = 4OCTET ; 参数个数，4 个字节
ARG = *OCTET ; 参数，字节数组形式
```

服务端：
```go
server := vex.NewServer()
server.RegisterHandler("test", func(ctx *vex.Context) {
	ctx.Write([]byte("Test!"))
})

err := server.ListenAndServe("tcp", ":5837")
if err != nil {
	panic(err)
}
```

客户端：

```go
client, err := vex.NewClient("tcp", "127.0.0.1:5837")
if err != nil {
	panic(err)
}
defer client.Close()

response, err := client.Do("test", [][]byte{
	[]byte("123"),
	[]byte("456"),
})

fmt.Println(string(response))
```

其实就是把 tcp 自定义通信做成一个模板，跟使用 http 框架类似，但是这个框架对于 tcp 通信来说性能不算很高，可以自行改进。

> R7-4700U，16GB 测试环境

```
BenchmarkServer-8          16258             71458 ns/op            8593 B/op         22 allocs/op
```