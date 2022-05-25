

# bytezero-go datachannel design.

## 服务端功能列表
[ ] 建立双向通道（client -> bytezero <- peer）
[ ] 缓存队列: ???
[ ] 身份验证
[ ] 数据传输： 创建Stream，CREATE -> OPEN
[ ] 关闭通道
[ ] API接口（http, jwt-token凭证）
- 获取创建通道参数AppId，AppKey，etc
- 获取通道流量信息
- 创建多监听者模式（解决单一tcp监听吞吐量）


## 客户端sdk(端到端传输)
[ ] 支持：java(android/linux/windows/macos)/golang/objectc/c++
[ ] 流程：
- 1、 建立通信双向通道（Channel）
- 1.1、客户端A发起请求（CHANNEL_CREATE, 设备deviceId，通道唯一标识sessionId等），Bytezero创建Channel，等待另一端连接或者10s超时
- 1.2、客户端B发起请求（CHANNEL_CREATE），Bytezero加入Channel，并通知A和B连接成功（CHANNEL_ACK）
- 2、 建立业务流通道（Stream）
- 2.1、客户端A发送创建流（STREAM_CREATE）消息，通过Bytezero转发给客户端B
- 2.2、客户端B接收到创建流消息，回调用户确认是否创建，并响应成功（STREAM_ACK)或拒绝（STREAM_ERROR)
- 2.3、客户端A接收到响应消息后通知用户成功（开始传输数据STREAM_DATA）或失败
- 3、 传输数据
- 3.1、二进制数据（文件数据等）
- 3.2、文本数据（控制协议等）
- 4、 关闭流通道
- 4.1、发起关闭流（STREAM_CLOSE）
- 5、 关闭通道（尽可能把数据发送出去）








