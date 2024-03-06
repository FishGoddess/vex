## ⛓ Vex

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/vex)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](./_icons/coverage.svg)
![Test](https://github.com/FishGoddess/vex/actions/workflows/test.yml/badge.svg)

**Vex** is a framework which uses tcp to exchange data.

[阅读中文版](./README.md)

### 🍃 Features

* Based on a tcp protocol, easy to use or customize
* Simple API design, connection pool supports
* Support client/server interceptors, easy to monitor and notify
* Signal monitor supports, shutdown gracefully
* Connection limit supports, and fast-failed supports
* Provided pack protocol, which is for simple data transmission protocol

_Check [HISTORY.md](./HISTORY.md) and [FUTURE.md](./FUTURE.md) to know about more information._

### 📃 Protocol

> The provided pack protocol defines a conception named packet no matter in request or response.

ABNF:

```abnf
PACKET = MAGIC TYPE DATASIZE DATA
MAGIC = 3OCTET ; 3Bytes, current is 0xC638B which is 811915
TYPE = OCTET ; 0x00-0xFF, begin from one, 255 at most
DATASIZE = 4OCTET ; 4bytes, 4GB at most
DATA = *OCTET ; Size is determined by DATASIZE
```

In human:

```
Packet:
magic     type    data_size    {data}
3byte     1byte     4byte      unknown
```

### 🔦 Examples

```bash
$ go get -u github.com/FishGoddess/vex
```

> We provide native and pack two ways to use: native is for customizing protocol and pack is a simple data transmission
> protocol.

Native client:

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

Native server:

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

Pack client:

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

Pack server:

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

_All examples can be found in [_examples](./_examples)._

### 🛠 Benchmarks

```bash
$ make bench
BenchmarkReadWrite-2      140317              8356 ns/op               0 B/op          0 allocs/op

$ make benchpack
BenchmarkPackReadWrite-2   61564             19650 ns/op            2080 B/op          6 allocs/op
```

| Protocol | Connections | rps          |
|----------|-------------|--------------|
| -        | 1           | &nbsp; 50231 |
| -        | 2           | 116790       |
| Pack     | 1           | &nbsp; 30852 |
| Pack     | 2           | &nbsp; 67453 |

_Packet size is 1KB._

_Environment: AMD EPYC 7K62, 2 Cores, 8GB RAM, linux._
