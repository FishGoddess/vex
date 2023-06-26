## ⛓ Vex

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/vex)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](./_icons/coverage.svg)
![Test](https://github.com/FishGoddess/vex/actions/workflows/test.yml/badge.svg)

**Vex** is a framework which uses tcp/udp to exchange data.

[阅读中文版](./README.md)

### 🍃 Features

* Based on a tcp/udp protocol, easy to use or customize
* Simple API design, connection pool supports
* Support client/server interceptors, easy to monitor and notify
* Signal monitor supports, shutdown gracefully
* Connection limit supports, provided several limit strategies.

_Check [HISTORY.md](./HISTORY.md) and [FUTURE.md](./FUTURE.md) to know about more information._

### 📄 Example

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
$ go test -v ./_examples/performance_test.go -bench=^BenchmarkServer$ -benchtime=1s
BenchmarkServer-16        136586              9063 ns/op            2080 B/op          6 allocs/op
```

> Packet size is 1KB.

_Environment: R7-5800X@3.8GHZ CPU, 32GB RAM, manjaro linux._

_Single connection: 10w requests spent 1.5s, result is **66876 rps**._

_Pool (16connections): 10w requests spent 359.9ms, result is **277859 rps**._
