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
* Server event callback supports, easy to monitor and notify.

_Check [HISTORY.md](./HISTORY.md) and [FUTURE.md](./FUTURE.md) to know about more information._

### 📃 Protocol

> All is packet including request and response.

ABNF：

```abnf
PACKET = HEADER BODY
HEADER = MAGIC VERSION TYPE BODYSIZE
BODY = *OCTET ; Size unknown, see BODYSIZE
MAGIC = 4OCTET ; 4Bytes, current is 0x755DD8C or 123067788
VERSION = OCTET ; 0x00-0xFF, begin from one, 255 at most
TYPE = OCTET ; 0x00-0xFF, begin from one, 255 at most
BODYSIZE = 4OCTET ; 4bytes, 4GB at most
```

In human:

```
Packet:
magic    version    type    body_size    {body}
4byte     1byte     1byte     4byte      unknown
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
	"fmt"

	"github.com/FishGoddess/vex"
)

func main() {
	server := vex.NewServer()
	server.RegisterPacketHandler(1, func(req []byte) (rsp []byte, err error) {
		fmt.Println(string(req))
		return []byte("server test"), nil
	})

	err := server.ListenAndServe("tcp", "127.0.0.1:5837")
	if err != nil {
		panic(err)
	}
}
```

* [client](./_examples/client.go)
* [server](./_examples/server.go)
* [pool](./_examples/pool.go)

_All examples can be found in [_examples](./_examples)._

### 🛠 Benchmarks

```bash
$ go test -v ./_examples/performance_test.go -bench=^BenchmarkServer$ -benchtime=1s
BenchmarkServer-16        187464              6758 ns/op              64 B/op          6 allocs/op
```

_Environment: R7-5800X@3.8GHZ CPU, 32GB RAM._

_Single connection: 10w requests spent 745.17ms, result is **134198 rps**, single spent 7.45us._

_Pool (16connections): 10w requests spent 277.06ms, result is **360933 rps**, single spent 2.77us._
