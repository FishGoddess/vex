## ⛓ Vex

> 基于 mit 协议开源

一个通用网络通信框架，可以作为参考案例进行交流学习哈哈。并发请求响应的支持需要比较复杂的协议设计，这个框架并不支持。

### 📃 协议描述

ABNF 描述：

```abnf
HEADER = VERSION CMD LENGTH ; 头部，主要是协议版本号，指令和参数长度
REQUESTBODY = *{ARGLENGTH ARG} ; 请求正文，主要是参数列表，每个参数都以长度开头
RESPONSEBODY = *OCTET ; 响应正文，字节数组
VERSION = OCTET ; 协议版本号，OCTET 是指 0x00-0xFF
CMD = OCTET ; 指令，0x00-0xFF，所以单个服务最多只能支持到 255 个指令
LENGTH = 4OCTET ; 参数个数，4 个字节
ARGLENGTH = 4OCTET ; 参数个数，4 个字节
ARG = *OCTET ; 参数，字节数组
```

人类语言描述：
```
请求：
version    command    argsLength    {argLength    arg}
 1byte      1byte       4byte          4byte    unknown

响应：
version    reply    bodyLength    {body}
 1byte     1byte      4byte      unknown
```

### ✒ 使用案例

服务端：
```go
server := vex.NewServer()
server.RegisterHandler(1, func(args [][]byte) (body []byte, err error) {
	return []byte("test"), nil
})

err := server.ListenAndServe("tcp", ":5837")
if err != nil {
	b.Fatal(err)
}
```

客户端：

```go
client, err := vex.NewClient("tcp", "127.0.0.1:5837")
if err != nil {
	b.Fatal(err)
}
defer client.Close()

response, err := client.Do(1, [][]byte{
	[]byte("123"), []byte("456"),
})
if err != nil {
	b.Fatal(err)
}

fmt.Println(string(response))
```

### 🛠 性能测试

其实就是把 tcp 自定义通信做成一个模板，跟使用 http 框架类似，但是这个框架对于 tcp 通信来说性能不算很高，可以自行改进。

> R7-4700U，16GB 测试环境

```
BenchmarkServer-8          53317             23556 ns/op             144 B/op         12 allocs/op
```