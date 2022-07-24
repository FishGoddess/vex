## ⛓ Vex

[![License](./_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Build](./_icons/build.svg)](./_icons/build.svg)
[![Coverage](./_icons/coverage.svg)](./_icons/coverage.svg)

**Vex** 是一个使用 tcp 通信和交换数据的框架。

[Read me in English](./README.en.md)

> 并发请求响应的支持需要比较复杂的协议设计，这个框架并不支持。

### 🍃 功能特性

* 基于 tcp 自定义通信协议，直接使用或二次开发都很简单
* 极简设计的 API，内置连接池，可以对性能进行调优
* 支持服务器事件回调机制，方便接入监控和告警
* 支持信号量监控机制，并支持平滑下线
* 支持服务器令牌桶连接数限制，并支持多种连接限制策略

_历史版本的特性请查看 [HISTORY.md](./HISTORY.md)。未来版本的新特性和计划请查看 [FUTURE.md](./FUTURE.md)。_

### 📃 协议描述

> 协议抽象出数据包的概念，不管是请求还是响应都视为一种数据包。

ABNF：

```abnf
PACKET = HEADER BODY ; 数据包
HEADER = MAGIC TYPE BODYSIZE ; 数据包头，主要是魔数，包类型以及包体大小
BODY = *OCTET ; 数据包体，大小未知，需要靠 BODYSIZE 来确认
MAGIC = 3OCTET ; 魔数，3 个字节表示，目前是 0xC638B，也就是 811915
TYPE = OCTET ; 数据包类型，0x00-0xFF，从 0 开始，最多 255 种数据包类型
BODYSIZE = 4OCTET ; 数据包体大小，4 个字节表示，最大是 4GB
```

人类语言描述：

```
数据包：
magic     type    body_size    {body}
3byte     1byte     4byte      unknown
```

### ✒ 使用案例

```bash
$ go get -u github.com/FishGoddess/vex
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

	rsp, err := client.Send(1, []byte("client test"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(rsp))
}
```

服务端：

```go
package main

import (
	"context"
	"fmt"

	"github.com/FishGoddess/vex"
)

func main() {
	server := vex.NewServer("tcp", "127.0.0.1:5837", vex.WithName("example"))
	server.RegisterPacketHandler(1, func(ctx context.Context, requestBody []byte) (responseBody []byte, err error) {
		addr, ok := vex.RemoteAddr(ctx)
		if !ok {
			fmt.Println(string(requestBody))
		} else {
			fmt.Println(string(requestBody), "from", addr)
		}
		return []byte("server test"), nil
	})

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
```

* [client](./_examples/client.go)
* [server](./_examples/server.go)
* [pool](./_examples/pool.go)
* [event](./_examples/event.go)

_所有的使用案例都在 [_examples](./_examples) 目录。_

### 🛠 性能测试

```bash
$ go test -v ./_examples/performance_test.go -bench=^BenchmarkServer$ -benchtime=1s
BenchmarkServer-16        136586              9063 ns/op            2080 B/op          6 allocs/op
```

> 数据包大小为 1KB。

_测试环境：R7-5800X@3.8GHZ CPU，32GB RAM，manjaro linux。_

_单连接：10w 个请求的执行耗时为 1.5s，结果为 **66876 rps**。_

_16个连接：10w 个请求的执行耗时为 359.9ms，结果为 **277859 rps**。_
