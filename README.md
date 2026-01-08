## ⛓ Vex

[![Go Doc](_icons/godoc.svg)](https://pkg.go.dev/github.com/FishGoddess/vex)
[![License](_icons/license.svg)](https://opensource.org/licenses/MIT)
[![Coverage](_icons/coverage.svg)](./_icons/coverage.svg)
![Test](https://github.com/FishGoddess/vex/actions/workflows/test.yml/badge.svg)

**Vex** 是一个使用 tcp 通信和传输数据的框架。

[Read me in English](./README.en.md)

### 🍃 功能特性

* 基于 tcp 自定义协议传输数据，使用简单
* 极简设计的 API，内置连接池，可以对性能进行调优
* 支持客户端、服务器引入拦截器，方便接入监控和告警
* 支持信号量监控机制和平滑下线
* 支持连接数限制，并支持超时中断

_历史版本的特性请查看 [HISTORY.md](./HISTORY.md)。未来版本的新特性和计划请查看 [FUTURE.md](./FUTURE.md)。_

### 📃 协议描述

ABNF：

```abnf
PACKET = MAGIC TYPE LENGTH SEQUENCE DATA
MAGIC = 3OCTET ; 魔数，目前是 0xC638B（811915）
TYPE = OCTET ; 类型，0x00-0xFF
LENGTH = 4OCTET ; 长度，最大 4GB
SEQUENCE = 8OCTET ; 序号，用来区分不同的数据包
DATA = *OCTET ; 数据，实际大小需要靠 LENGTH 来确认
```

人话：

```
数据包：
magic     type     length    sequence     {data}
3byte     1byte     4byte     8byte       unknown
```

_你会发现协议没有版本号的字段，其实是我们选择将版本号融入到类型字段中，因为每个版本可能对应的类型都不一样。_

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
```

### 👥 贡献者

如果您觉得 goes 缺少您需要的功能，请不要犹豫，马上参与进来，发起一个 _**issue**_。
