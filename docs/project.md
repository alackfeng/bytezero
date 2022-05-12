

# bytezero-go - 中转数据传输通道网络模型.

## 命令行参数&读取配置文件.
`````
#### 设置proxy
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

#### 构建命令行参数工程.
cd D:\taurus\bitdisk\code\bytezero\
go get -u github.com/spf13/cobra/cobra
cobra init bytezero --pkg-name github.com/alackfeng/bytezero

#### 添加子服务.
cd D:\taurus\bitdisk\code\bytezero\bytezero
cobra add server -t github.com/alackfeng/bytezero
cobra add client -t github.com/alackfeng/bytezero
cobra add tool -t github.com/alackfeng/bytezero

#### 读取配置.


#### 运行.
go run .\main.go server

go build -o bin/bytezero.exe -v main.go
.\bin\bytezero.exe server


`````

## 使用go mod管理工程.
`````

cd D:\taurus\bitdisk\code\bytezero\bytezero

#### 初始化go.mod.
go mod init github.com/alackfeng/bytezero

#### 检查依赖源码包
go mod tidy -v

#### 
go build main.go

#### 列出依赖
go list -m all
go list -m -versions golang.org/x/text

#### 更新依赖
go get golang.org/x/text
go get golang.org/x/text@v1.3.1

####
go mod verify

####
go mod vendor -v

`````


## 



