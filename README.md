## ⛓ Vex

[![License](./_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Build](./_icons/build.svg)](./_icons/build.svg)
[![Coverage](./_icons/coverage.svg)](./_icons/coverage.svg)

把 tcp 自定义通信协议做成一个模板，跟使用 http 框架类似，只不过他的性能非常强悍，也算是一个通用网络通信框架，大家可以作为参考案例进行交流学习哈哈。

> 并发请求响应的支持需要比较复杂的协议设计，这个框架并不支持。

### 📃 协议描述

ABNF 描述请求：

```abnf
REQUEST = HEADER BODY ; 请求
HEADER = VERSION COMMAND ARGSLENGTH ; 请求头，主要是版本号，命令以及参数个数
BODY = *{ARGLENGTH ARG} ; 请求体，主要是参数，*{} 表示可能 {} 里面的东西可能没有，也可能有多个
VERSION = OCTET ; 版本号，0x00-0xFF，一般从 1 开始，也就是最多 255 个版本号
COMMAND = OCTET ; 命令，0x00-0xFF，一般从 1 开始，也就是最多 255 个命令
ARGSLENGTH = 4OCTET ; 参数个数，4 个字节表示，也就是最多 uint32 个参数
ARGLENGTH = 4OCTET ; 参数长度，4 个字节表示，也就是最长是 uint32 个字节
ARG = *OCTET ; 参数内容，长度未知，需要靠 ARGLENGTH 明确
```

ABNF 描述响应：

```abnf
RESPONSE = HEADER BODY ; 响应
HEADER = VERSION REPLY BODYLENGTH ; 响应头，主要是版本号，命令以及参数个数
BODY = *OCTET ; 响应体，长度未知，需要靠 BODYLENGTH 明确
VERSION = OCTET ; 版本号，0x00-0xFF，一般从 1 开始，也就是最多 255 个版本号
REPLY = OCTET ; 命令，0x00-0xFF，一般从 1 开始，也就是最多 255 种答复含义
BODYLENGTH = 4OCTET ; 参数长度，4 个字节表示，也就是最长是 uint32 个字节
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
package main

import "github.com/FishGoddess/vex"

func main() {

	server := vex.NewServer()
	server.RegisterHandler(1, func(args [][]byte) (body []byte, err error) {
		return []byte("test"), nil
	})

	err := server.ListenAndServe("tcp", ":5837")
	if err != nil {
		panic(err)
	}
}
```

客户端：

```go
package main

import (
	"fmt"

	"github.com/FishGoddess/vex"
)

func main() {

	client, err := vex.NewClient("tcp", "127.0.0.1:5837")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	response, err := client.Do(1, [][]byte{
		[]byte("123"), []byte("456"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(response))
}
```

### 🛠 性能测试

```bash
$ go test -v ./_examples/performance_test.go -run=^$ -bench=^BenchmarkServer$ -benchtime=1s
BenchmarkServer-16        112126             12759 ns/op             144 B/op         11 allocs/op
```

_测试环境：R7-5800X@3.8GHZ CPU，32GB RAM。_

_单连接：10000 个命令的执行耗时为 114.02ms，得到的 rps 为 **87704 w/s**，单命令耗时 11.402 us。_

_连接池（64个连接）：并发 10000 个命令的执行耗时为 20.04ms，得到的 rps 为 **499001 w/s**，单命令耗时 2.004 us。_
