## ⛓ Vex

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/vex)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](./_icons/coverage.svg)
![Test](https://github.com/FishGoddess/vex/actions/workflows/test.yml/badge.svg)

**Vex** 是一个使用 tcp 通信和传输数据的框架。

[Read me in English](./README.en.md)

### 🍃 功能特性

* 基于 tcp 自定义协议传输数据，极简 API 设计
* 支持信号量监控机制和平滑下线
* 支持连接数限制，并支持超时中断（敬请期待）
* 支持客户端、服务器两种拦截器，方便监控（敬请期待）
* 内置连接池，可以对性能进行调优（敬请期待）

_历史版本的特性请查看 [HISTORY.md](./HISTORY.md)。未来版本的新特性和计划请查看 [FUTURE.md](./FUTURE.md)。_

### 📃 协议描述

ABNF：

```abnf
PACKET = ID MAGIC FLAGS LENGTH DATA
ID = 8OCTET ; 编号，用来区分不同的数据包
MAGIC = 4OCTET ; 魔数，目前是 1997811915
FLAGS = 8OCTET ; 标志位，比如是否为错误包
LENGTH = 4OCTET ; 长度，最大 4GB
DATA = *OCTET ; 数据，需要靠 LENGTH 来确认
```

人话：

```
数据包：
id       magic     flags     length     {data}
8byte    4byte     8byte     4byte      unknown
```

_你会发现协议没有版本号的字段，其实是我们选择将版本号融入到魔数字段中，所以每个版本可能对应的魔数不一样。_

### 🔦 使用案例

```bash
$ go get -u github.com/FishGoddess/vex
```

客户端：

```go
package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/FishGoddess/vex"
)

func main() {
	client, err := vex.NewClient("127.0.0.1:9876")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	ctx := context.Background()
	for i := range 10 {
		data := []byte(strconv.Itoa(i))
		fmt.Printf("client send: %s\n", data)

		data, err = client.Send(ctx, data)
		if err != nil {
			panic(err)
		}

		fmt.Printf("client receive: %s\n", data)
		time.Sleep(100 * time.Millisecond)
	}
}
```

服务端：

```go
package main

import (
	"context"

	"github.com/FishGoddess/vex"
)

type EchoHandler struct{}

func (EchoHandler) Handle(ctx context.Context, data []byte) ([]byte, error) {
	return data, nil
}

func main() {
	server := vex.NewServer("127.0.0.1:9876", EchoHandler{})
	defer server.Close()

	if err := server.Serve(); err != nil {
		panic(err)
	}
}
```

_所有的使用案例都在 [_examples](./_examples) 目录。_

### 🛠 性能测试

```bash
$ make bench
```

```bash
goos: linux
goarch: amd64
cpu: Intel(R) Xeon(R) CPU E5-26xx v4

BenchmarkPacket-2          29292             38818 ns/op            4600 B/op          9 allocs/op
```

> 测试文件：[_examples/packet_test.go](./_examples/packet_test.go)。

### 👥 贡献者

如果您觉得 goes 缺少您需要的功能，请不要犹豫，马上参与进来，发起一个 _**issue**_。
