
# 测试命令.


## 服务端程序.

go run .\main.go server


## 客户端程序.

go run .\main.go client -t 127.0.0.1:7788 -s bytezero-session-id-0 -d device0
go run .\main.go client -t 127.0.0.1:7788 -s bytezero-session-id-0 -d device1
