## ⛓ Vex

[![License](./_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Build](./_icons/build.svg)](./_icons/build.svg)
[![Coverage](./_icons/coverage.svg)](./_icons/coverage.svg)

**Vex** is a framework which uses tcp to send packets and exchange data through two processes.

[阅读中文版](./README.md)

> Concurrent protocol is too complex and vex doesn't support.

### 🥇 Features

* Based on a customized tcp protocol, easy to use and develop
* Simple API design, client pool supports
* Server event callback supports, easy to monitor and notify
* Signal supports, shutdown server gracefully
* Server connection token supports, and three limit strategies supports.

_Check [HISTORY.md](./HISTORY.md) and [FUTURE.md](./FUTURE.md) to know about more information._

### 📃 Protocol

> All is packet including request and response.

ABNF：

```abnf
PACKET = HEADER BODY
HEADER = MAGIC TYPE BODYSIZE
BODY = *OCTET ; Size unknown, see BODYSIZE
MAGIC = 3OCTET ; 3Bytes, current is 0xC638B
TYPE = OCTET ; 0x00-0xFF, begin from one, 255 at most
BODYSIZE = 4OCTET ; 4bytes, 4GB at most
```

In human:

```
Packet:
magic     type    body_size    {body}
3byte     1byte     4byte      unknown
```

### ✒ Example

```bash
$ go get -u github.com/FishGoddess/vex
```

Client:

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

Server:

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

_All examples can be found in [_examples](./_examples)._

### 🛠 Benchmarks

```bash
$ go test -v ./_examples/performance_test.go -bench=^BenchmarkServer$ -benchtime=1s
BenchmarkServer-16        136586              9063 ns/op            2080 B/op          6 allocs/op
```

> Packet size is 1KB.

_Environment: R7-5800X@3.8GHZ CPU, 32GB RAM, manjaro linux._

_Single connection: 10w requests spent 1.5s, result is **66876 rps**._

_Pool (16connections): 10w requests spent 359.9ms, result is **277859 rps**._
