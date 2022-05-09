

# bytezero-go datachannel design.

## 功能列表
[ ] 申请临时通道地址（jwt-token凭证）
- 监听端口： http(s)\udp\tcp(ssl)
- 分配通道：根据客户端传入的参数（设备id，通道唯一标识sessionid等）分配通道（临时端口，token凭证等）
- 信任问题：AppID，AppKey
[ ] 建立双向通道（client -> bytezero <- peer）
- client申请连接：client连接临时端口，等待peer或者超时
- peer响应连接： peer连接临时端口，等待client或超时，并通过bytezero转发client连接请求
- 双向确认： peer接收连接请求后，回ACK给client。bytezero转发ACK响应，client变更状态Connected，可传递消息
- 身份验证： 携带token凭证、设备id等，验证身份
- 定时心跳: 防止断开
- 控制协议：可以控制速率bps等
[ ] 数据传输（保证时序）
- 传输数据： 创建任务标签，CREATE -> ACK
[ ] 关闭通道
- 主动关闭（保证所有数据被接收后关闭）
- 被动关闭

## 重点问题
- 消息时序
- 缓存队列


## 客户端sdk(端到端传输)
- 支持：java(android/linux/windows/macos)/golang/objectc






