## ⛓ Vex

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/vex)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](./_icons/coverage.svg)
![Test](https://github.com/FishGoddess/vex/actions/workflows/test.yml/badge.svg)

**Vex** is a framework using tcp to transfer data.

[阅读中文版](./README.md)

### 🍃 Features

* Based on a tcp protocol, easy to use
* Simple API design, connection pool supports
* Support client/server interceptors, easy to monitor and notify
* Signal monitor supports, shutdown gracefully
* Connection limit supports, and fast-failed supports

_Check [HISTORY.md](./HISTORY.md) and [FUTURE.md](./FUTURE.md) to know about more information._

### 📃 Protocol

ABNF:

```abnf
PACKET = ID MAGIC FLAGS LENGTH DATA
ID = 8OCTET ; Identify different packets
MAGIC = 3OCTET ; value is 1997811915
FLAGS = 8OCTET ; Set some flags of packet
LENGTH = 4OCTET ; 4GB at most
DATA = *OCTET ; Determined by LENGTH
```

In human:

```
Packet:
id       magic     flags     length     {data}
8byte    3byte     1byte     4byte      unknown
```

_The version of protocol is in magic because we think different versions may have different magics._

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
```

### 👥 Contributing

If you find that something is not working as expected please open an _**issue**_.
