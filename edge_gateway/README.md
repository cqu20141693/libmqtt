## 打包
### mac
1. m1
go build -o gateway gateway.go
### windows

## 功能扩展
### edge_gateway mqtt cloud
参考 GPlatformClient

### 网络异常重连
参考 NetHandlerWithReconnect

### 报文发送失败重试&&一直失败缓存
// 在//core:send 注释处处理err
当err 为连接不可用时
将消息丢入重试队列中； 进行重试3次
重试3次失败则丢入缓存队列中。


### 报文ack 失败重发机制
参考 idGenerator 通过next 方法生成id 进行发送，通过free 方法释放id;
过程中说明协议层已经发送完成。 构建本地id session;
