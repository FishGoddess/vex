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

### 📃 协议描述

> 协议抽象出数据包的概念，不管是请求还是响应都视为一种数据包。

ABNF：

```abnf
PACKET = HEADER BODY ; 数据包
HEADER = MAGIC TYPE BODYSIZE ; 数据包头，主要是魔数，包类型以及包体大小
BODY = *OCTET ; 数据包体，大小未知，需要靠 BODYSIZE 来确认
MAGIC = 3OCTET ; 魔数，3 个字节表示，目前是 0xC638B，也就是 811915
TYPE = OCTET ; 数据包类型，0x00-0xFF，从 0 开始，最多 255 种数据包类型
BODYSIZE = 4OCTET ; 数据包体大小，4 个字节表示，最大是 4GB
```

人类语言描述：

```
数据包：
magic     type    body_size    {body}
3byte     1byte     4byte      unknown
```

### 🔦 使用案例

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
$ make bench
BenchmarkReadWrite-16             183592              6603 ns/op               0 B/op          0 allocs/op
```

> 数据包大小为 1KB。

_测试环境：R7-5800X@3.8GHZ CPU, 32GB RAM, deepin linux。_

_单连接：10w 个请求的执行耗时为 1.26s，结果为 **78958 rps**。_

_16个连接：10w 个请求的执行耗时为 393.08ms，结果为 **254400 rps**。_
