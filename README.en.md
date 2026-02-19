## ⛓ Vex

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/vex)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](./_icons/coverage.svg)
![Test](https://github.com/FishGoddess/vex/actions/workflows/test.yml/badge.svg)

**Vex** is a framework using tcp to transfer data.

[阅读中文版](./README.md)

### 🍃 Features

* Based on a vex tcp protocol, simple API design
* Signal monitor supports, shutdown gracefully
* Connection limit supports, and timeout supports (Coming Soon)
* Support client/server interceptors, easy to observe (Coming Soon)
* Connection pool supports (Coming Soon)

_Check [HISTORY.md](./HISTORY.md) and [FUTURE.md](./FUTURE.md) to know about more information._

### 📃 Protocol

ABNF:

```abnf
PACKET = ID MAGIC FLAGS LENGTH DATA
ID = 8OCTET ; Identify different packets
MAGIC = 4OCTET ; value is 1997811915
FLAGS = 8OCTET ; Set some flags of packet
LENGTH = 4OCTET ; 4GB at most
DATA = *OCTET ; Determined by LENGTH
```

In human:

```
Packet:
id       magic     flags     length     {data}
8byte    4byte     8byte     4byte      unknown
```

_The version of protocol is in magic because we think different versions may have different magics._

### 🔦 Examples

```bash
$ go get -u github.com/FishGoddess/vex
```

Client:

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

Server:

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

_All examples can be found in [_examples](./_examples)._

### 🛠 Benchmarks

```bash
$ make bench
```

```bash
goos: linux
goarch: amd64
cpu: Intel(R) Xeon(R) CPU E5-26xx v4

BenchmarkPacket-2          48885             25712 ns/op            4600 B/op          9 allocs/op
BenchmarkPacketPool-2      58665             21461 ns/op            4601 B/op          9 allocs/op
```

> Benchmark: [_examples/packet_test.go](./_examples/packet_test.go).

> Pool benchmark uses 2 clients and the network card is the bottleneck.

### 👥 Contributing

If you find that something is not working as expected please open an _**issue**_.
