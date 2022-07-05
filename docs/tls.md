
## TLS/SSL加密传输



#### https Client

#### https Server


#### 生成自签名self-signed证书 根CA
`````
mkdir -p scripts/certs
cd scripts/certs

- 1 生成CA自己的私钥 rootCA.key
openssl genrsa -out rootCA.key 2048
- 2 根据CA自己的私钥生成自签发的数字证书，该证书里包含CA自己的公钥
openssl req -x509 -new -nodes -key rootCA.key -subj "/CN=*.bytezero.org" -days 5000 -out rootCA.pem

- 3. 生成服务端的私钥和数字证书（由自CA签发）
# 生成服务端私钥
openssl genrsa -out server.key 2048
# 生成Certificate Sign Request，CSR，证书签名请求
openssl req -new -key server.key -subj "/CN=localhost" -out server.csr
# 自CA用自己的CA私钥对服务端提交的csr进行签名处理，得到服务端的数字证书server.crt
openssl x509 -req -in server.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out server.crt -days 5000

- 4. 将自CA的数字证书同客户端一并发布，用于客户端对服务端的数字证书进行校验
mkdir -p /client
cp rootCA.pem ../certs/client/caroot.crt # or caroot.pem
- 5. 将服务端的数字证书和私钥同服务端一并发布
mkdir -p server/
mv server.crt ../certs/server/server.crt 
mv server.key ../certs/server/server.key

`````

#### 对客户端的证书进行校验(双向证书校验）
`````
openssl genrsa -out client.key 2048
openssl req -new -key client.key -subj "/CN=bytezeroclient_cn" -out client.csr
openssl x509 -req -in client.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out client.crt -days 5000

mv client.crt ../certs/client/client.crt 
mv client.key ../certs/client/client.key

`````

#### openssl 生产证书
`````
- 1. 生成私钥(Private Key) .key 
openssl genrsa -out server.key 2048

- 2. 生成证书(Certificate) .crt
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

-----
Country Name (2 letter code) []:cn
State or Province Name (full name) []:Guangdong
Locality Name (eg, city) []:Guangzhou
Organization Name (eg, company) []:bitdisk.io
Organizational Unit Name (eg, section) []:bitdisk.io
Common Name (eg, fully qualified host name) []:bitdisk.io
Email Address []:info@bitdisk.io 
➜  docs git:(master) ✗ 


- 3. 生成pem和key
openssl req -new -nodes -x509 -out server.pem -keyout server.key -days 3650 -subj "/C=CN/ST=Guangdong/L=Guangzhou/O=bitdisk.io/OU=bitdisk/CN=localhost/emailAddress=admin@bitdisk.io"

- 4. 查看PEM格式证书的信息
openssl x509 -in server.pem -text -noout

`````

