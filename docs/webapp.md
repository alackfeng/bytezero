# nuxtjs h5 app

## https://www.nuxtjs.cn/guide/installation

## https://typescript.nuxtjs.org/cookbook/plugins

```
npx create-nuxt-app webapp

#### To get started:

cd nuxtjs-app
npm run dev

#### To build & start for production:

cd nuxtjs-app
npm run build
npm run start

#### To test:

cd nuxtjs-app
npm run test

```

## 优化.

```

#### stylelint.
npm install -D stylelint@latest stylelint-config-standard@latest stylelint-order@latest stylelint-config-recess-order@latest postcss-less@latest
npm install -D postcss-html stylelint-config-recommended-vue
npm install -D @nuxtjs/stylelint-module

#### i18n.
npm install @nuxtjs/i18n

#### vconsole
npm install vconsole

#### 图表.
npm install echarts

```

## 依赖.

```
npm install @vant/touch-emulator
npm install git+ssh://git@gitlab.cume.cc:bytezeroc/peerjs.git
npm install pako
npm install streamsaver @types/streamsaver
npm install @nuxtjs/device
npm install @nuxtjs/i18n
npm install less less-loader@7

npm install express
npm install qrcode

```

## 实时带宽评估算法调研.
`````
- 1. 获取网卡支持最大带宽100Mb/s 1000Mb/s, 及当前实时速率
可通过调用系统命令ethtool获取到当前网卡的最大支持带宽
- 2. 向服务发起实时带宽测试（上传/下载用例), 获取当前可支持的带宽容量
- 3. 桶算法进行带宽分配
- 4. 例如：
- 4.1 如当前带宽100Mb/s, 预留20%，当前可支持80Mb/s（当前带宽值是根据最大带宽及实施速率得出最小值，或者可考虑历史评估情况）；
- 4.2 每秒向桶内插入80Mb的容量，最大80Mb;
- 4.3 每个任务在数据传输过程中先向桶取可用容量，没有就等待1s在取
- 4.4 新建立的任务发现没有容量直接返回繁忙.


`````

## 带宽评估算法.
## 参考：https://developer.aliyun.com/article/892402
## 参考：https://zhuanlan.zhihu.com/p/515968239
`````
1. 传统的基于策略的算法GCC
- 延时估计和丢包相结合的拥塞控制算法。基于datachannel的webrtc可以在业务传输协议上增加timestamp，接收端计算延迟情况对带宽进行评估，并反馈给发送端。
这样确实可以解决每路传输线路上的带宽，也就解决了总体带宽被拥塞的问题。
- 另外新连接与带宽速率是有线性关系。随着连接数增加，带宽速率降低.
- 总带宽的本身的意义只是说明可提供的最大带宽容量，对于复杂的家庭网络意义不是很大。动态评估每一路传输的可用性更重要些吧。
- 历史数据统计也是很有意义的，可以对某些指标间的关系进行评估，以便更好的调整参数

2. 基于学习的模型 
- 不太实用于弱设备

`````

## 网络性能指标.
`````
- 带宽 bps
100Mbps = 12.5MB/s

- 时延
ping: icmp报文
处理时延、排队时延、发送时延、传播时延

- 抖动
网络抖动是指最大延迟与最小延迟的时间差

- 丢包

----
probe探测算法

---- ​WebRTC 拥塞控制 | 计算包组时间差 - InterArrival
https://developer.aliyun.com/article/781505?spm=a2c6h.14164896.0.0.3f5f7ff6AX2Chx

---- WebRTC 拥塞控制 | Trendline 滤波器
https://developer.aliyun.com/article/781509?spm=a2c6h.14164896.0.0.3f5f7ff6AX2Chx

---- WebRTC 拥塞控制 | 网络带宽过载检测
https://developer.aliyun.com/article/781511?spm=a2c6h.14164896.0.0.3f5f7ff6AX2Chx3


---- WebRTC 拥塞控制 | AIMD 码率控制
https://developer.aliyun.com/article/781538?spm=a2c6h.14164896.0.0.3f5f7ff6AX2Chx


---- 测试工具：
- 1. iperf3
- 2. https://www.speedtest.cn/


总带宽
实时总带宽


`````



