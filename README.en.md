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
* Connection limit supports, provided several limit strategies.

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

### 🔦 Examples

```bash
$ go get -u github.com/FishGoddess/vex
```

Client:

```go
```

Server:

```go
```

_All examples can be found in [_examples](./_examples)._

### 🛠 Benchmarks

```bash
$ make bench
BenchmarkReadWrite-16             183592              6603 ns/op               0 B/op          0 allocs/op
```

> Packet size is 1KB.

_Environment: R7-5800X@3.8GHZ CPU, 32GB RAM, deepin linux._

_Single connection: 10w requests spent 1.26s, result is **78958 rps**._

_Pool (16connections): 10w requests spent 393.08ms, result is **254400 rps**._
