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
* 支持连接数限制，并支持多种限制策略

_历史版本的特性请查看 [HISTORY.md](./HISTORY.md)。未来版本的新特性和计划请查看 [FUTURE.md](./FUTURE.md)。_

### 📄 使用案例

```bash
$ go get -u github.com/FishGoddess/vex
```

客户端：

```go
```

服务端：

```go
```

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
