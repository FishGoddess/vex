## ⛓ Vex

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/vex)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](./_icons/coverage.svg)
![Test](https://github.com/FishGoddess/vex/actions/workflows/test.yml/badge.svg)

**Vex** 是一个使用 tcp 通信和传输数据的框架。

[Read me in English](./README.en.md)

### 🍃 功能特性

* 基于 tcp 传输数据，直接使用或二次开发都很简单
* 极简设计的 API，内置连接池，可以对性能进行调优
* 支持客户端、服务器引入拦截器，方便接入监控和告警
* 支持信号量监控机制和平滑下线
* 支持连接数限制，并支持超时中断
* 自带 pack 数据传输协议，用于简单的数据传输场景

_历史版本的特性请查看 [HISTORY.md](./HISTORY.md)。未来版本的新特性和计划请查看 [FUTURE.md](./FUTURE.md)。_

### 📃 协议描述

> 自带的 pack 数据传输协议抽象出了一个数据包的概念，不管是请求还是响应都视为一种数据包。

ABNF：

```abnf
PACKET = MAGIC TYPE DATASIZE DATA ; 数据包
MAGIC = 3OCTET ; 魔数，3 个字节表示，目前是 0xC638B，也就是 811915
TYPE = OCTET ; 数据包类型，0x00-0xFF，从 0 开始，最多 255 种数据包类型
DATASIZE = 4OCTET ; 数据包的数据大小，4 个字节表示，最大是 4GB
DATA = *OCTET ; 数据包的数据，大小未知，需要靠 DATASIZE 来确认
```

人类语言描述：

```
数据包：
magic     type    data_size    {data}
3byte     1byte     4byte      unknown
```

### 🔦 使用案例

```bash
$ go get -u github.com/FishGoddess/vex
```

> 我们提供了原生和 pack 两种使用方式，其中原生可以自定义协议，随意读写操作数据，用于二次开发，而 pack
> 则是自带的数据包传输协议，用于简单的数据传输场景。

原生客户端：

```go
package main

import (
	"fmt"

	"github.com/FishGoddess/vex"
)

func main() {
	client, err := vex.NewClient("127.0.0.1:6789")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	msg := []byte("hello")
	if _, err := client.Write(msg); err != nil {
		panic(err)
	}

	var buf [1024]byte
	n, err := client.Read(buf[:])
	if err != nil {
		panic(err)
	}

	fmt.Println("Received:", string(buf[:n]))
}
```

原生服务端：

```go
package main

import (
	"fmt"
	"io"

	"github.com/FishGoddess/vex"
)

func handle(ctx *vex.Context) {
	var buf [1024]byte
	for {
		n, err := ctx.Read(buf[:])
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		fmt.Println("Received:", string(buf[:n]))

		if _, err = ctx.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
}

func main() {
	// Create a server listening on 127.0.0.1:6789 and set a handle function to it.
	// Also, we can give it a name like "echo" so we can see it in logs.
	server := vex.NewServer("127.0.0.1:6789", handle, vex.WithName("echo"))

	// Use Serve() to begin serving.
	// Press ctrl+c/control+c to close the server.
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
```

Pack 客户端：

```go
package main

import (
	"fmt"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pack"
)

func main() {
	client, err := vex.NewClient("127.0.0.1:6789")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	// Use Send method to send a packet to server and receive a packet from server.
	// Try to change 'hello' to 'error' and see what happens.
	packet, err := pack.Send(client, 1, []byte("hello"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(packet))
}
```

Pack 服务端：

```go
package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/pack"
)

func newRouter() *pack.Router {
	router := pack.NewRouter()

	// Use Register method to register your handler for some packets.
	router.Register(1, func(ctx context.Context, packetType pack.PacketType, requestPacket []byte) (responsePacket []byte, err error) {
		msg := string(requestPacket)
		fmt.Println(msg)

		if msg == "error" {
			return nil, errors.New(msg)
		} else {
			return requestPacket, nil
		}
	})

	return router
}

func main() {
	// Create a router for packets.
	router := newRouter()

	// Create a server listening on 127.0.0.1:6789 and set a handle function to it.
	server := vex.NewServer("127.0.0.1:6789", router.Handle, vex.WithName("pack"))

	// Use Serve() to begin serving.
	// Press ctrl+c/control+c to close the server.
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
```

_所有的使用案例都在 [_examples](./_examples) 目录。_

### 🛠 性能测试

```bash
$ make bench
BenchmarkReadWrite-16             172698              6795 ns/op               0 B/op          0 allocs/op
BenchmarkPackReadWrite-16          76129             16057 ns/op            2080 B/op          6 allocs/op
```

| 协议   | 连接数      | rps          |
|------|----------|--------------|
| -    | &nbsp; 1 | &nbsp; 77128 |
| -    | 16       | 256088       |
| Pack | &nbsp; 1 | &nbsp; 49796 |
| Pack | 16       | 200490       |

_数据包大小为 1KB。_

_测试环境：R7-5800X@3.8GHZ CPU, 32GB RAM, deepin linux。_
