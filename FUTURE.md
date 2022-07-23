## ✒ 未来版本的新特性 (Features in future versions)

### v0.3.x

* [x] 修改协议，去除 version，将 magic/type/body_size 混到一个 int64 中，使用位操作处理数值
* [x] PacketHandler 加入 context 并默认传递一些数据，比如客户端的地址
* [x] 提供一个 context 工具包，获取常用数据
* [x] 给 EventHandler 加入 context 参数，使用 context 传递数据而不是用 source 机制
* [x] 独立 pool 包的 config 结构
* [x] 考虑连接池的拒绝策略必要性，只保留阻塞和失败两种
* [x] 增加连接池 wait get 的请求数量
* [x] 修复性能压测超过 32KB 就出问题的 bug（io.Reader 的 Read 方法不保证读取满 slice 长度，要使用 io.ReadFull 才行）
* [x] 增加关闭服务器超时机制，防止连接过多导致关闭阻塞卡死
* [ ] 完善连接池的实现，加入 context 超时（发现在 select 中增加一个 case 会导致性能急剧下降。。。原因是 runtime.selectgo 方法）
* [ ] 给 Server 加入令牌桶模式的连接数控制，完善拒绝策略

### v0.2.x

* [x] 完善文档注释
* [x] 完善网络通信协议的设计
* [x] 加入连接池状况查询入口
* [x] 完善连接池的实现，支持基础数量限制
* [x] 考虑将 Client 做成接口
* [x] 加入 Server 的事件回调机制
* [x] 加入 signal 信号监听，引入平滑下线机制
* [x] 给 Client 和 Server 加入 option 机制
* [x] 抽象事件处理器，配置默认事件处理器

### v0.2.0-alpha

* [x] 实现最简单的网络通信功能
* [x] 加入简单的连接池实现