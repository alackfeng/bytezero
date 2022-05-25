package client

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/alackfeng/bytezero/bytezero/protocol"
	"github.com/alackfeng/bytezero/cores/utils"
)


const networkTcp = "tcp"
const waitTimeoutForAck = 10000 // ms.

// AppsChannel -
type AppsChannel struct {
    app Client
    *net.TCPConn
    address string
    sessionId string
    cid protocol.ChannelId
    ack chan protocol.ErrCode
    state protocol.ChannelState

    l sync.Mutex
    m map[protocol.StreamId]*AppsStream
    w map[protocol.StreamId]chan protocol.ErrCode

    // stats.
    sendStat utils.StatBandwidth
    recvStat utils.StatBandwidth
}
var _ ChannelHandle = (*AppsChannel)(nil)
var _ ChannelSender = (*AppsChannel)(nil)

// NewAppsChannel -
func NewAppsChannel(app Client) *AppsChannel {
    return &AppsChannel{
        app: app,
        m: make(map[protocol.StreamId]*AppsStream),
        w: make(map[protocol.StreamId]chan protocol.ErrCode),
    }
}

// Address -
func (a *AppsChannel) Address() string {
    return a.address
}

// Send -
func (a *AppsChannel) Send(buf []byte) error {
    n, err := a.Write(buf)
    if err != nil {
        return err
    }
    if n != len(buf) {
        return protocol.ErrBufferNotAllSent
    }
    return nil
}

// Start -
func (a *AppsChannel) Start(address, sessionId string) error {
    // 1. 建立Tcp连接.
    a.address = address
    if err := a.dial(); err != nil {
        return err
    }
    // 2. 启动消息接收.
    go a.handleRecevicer()

    // 3. 建立通信双向通道（Channel），连接对端设备.
    if err := a.channelCreate(sessionId); err != nil {
        return err
    }
    return a.wait(a.ack, waitTimeoutForAck, "channel", 0)
}

// Stop -
func (a *AppsChannel) Stop() error {
    return a.channelClose()
}

// Online -
func (a *AppsChannel) Online() bool {
    return a.state == protocol.ChannelStateOpen
}

// dial -
func (a *AppsChannel) dial() error {
    tcpAddr, err := net.ResolveTCPAddr(networkTcp, a.address)
    if err != nil {
        return err
    }
    tcpConn, err := net.DialTCP(networkTcp, nil, tcpAddr)
    if err != nil {
        return err
    }
    a.TCPConn = tcpConn
    return nil
}

// wait -
func (a *AppsChannel) wait(ack chan protocol.ErrCode, ms int, label string, sid protocol.StreamId) error {
    duration := time.Duration(ms) * time.Millisecond
    ticker := time.NewTicker(duration)
    select {
    case code := <- ack:
        if a.Online() {
            return nil
        }
        if label == "channel" {
            return fmt.Errorf("Channel#%s, errcode.%d", a.sessionId, code)
        } else {
            return fmt.Errorf("Stream#%d, errcode.%d", sid, code)
        }
    case <- ticker.C:
        if label == "channel" {
            return fmt.Errorf("Channel#%s timeout", a.sessionId)
        } else {
            return fmt.Errorf("Stream#%d timeout", sid)
        }
    }
}

// handleRecevicer -
func (a *AppsChannel) handleRecevicer() {
    if a.address == "" {
        return
    }

    recvBufferLen := a.app.MaxRecvBufferLen()
    buffer := make([]byte, recvBufferLen)
    currTime := time.Now()
    for {
        n, err := a.Read(buffer)
        if err != nil {
            fmt.Printf("AppsChannel.handleRecevicer error.%v.\n", err.Error())
            break
        }
        if a.recvStat.Bytes == 0 {
            a.recvStat.Begin()
            fmt.Printf("AppsChannel.handleRecevicer - begin. recv begin %v.\n", a.recvStat.Info())
        }
        a.recvStat.Inc(int64(n))
        if n != recvBufferLen {
            fmt.Printf("AppsChannel.handleRecevicer recv buffer len %d not equal send buffer, real %d.\n", recvBufferLen, n)
        }
        if time.Now().Sub(currTime).Milliseconds() > 1000 {
            currTime = time.Now()
            fmt.Printf("AppsChannel.handleRecevicer recv - count %d, bps %s. send bps %s\n", a.recvStat.Count, utils.ByteSizeFormat(a.recvStat.Bps1s()), utils.ByteSizeFormat(a.sendStat.Bps1s()))
        }

        // 处理接收到消息.
        commonPt := &protocol.CommonPt{}
        if err := commonPt.Unmarshal(buffer[0:n]); err != nil {
            fmt.Printf("AppsChannel.handleRecevicer - Unmarshal Buffer error.%v.\n", err.Error())
            continue
        }
        if err := a.handlePt(commonPt); err != nil {
            fmt.Printf("AppsChannel.handleRecevicer - handlePt error.%v.\n", err.Error())
            continue
        }

    }
    a.recvStat.End()
    fmt.Printf("AppsChannel.handleRecevicer - end... %v.\n", a.recvStat.InfoAll())
    // a.done <- true
}

// handlePt -
func (a *AppsChannel) handlePt(commonPt *protocol.CommonPt) error {
    switch commonPt.Type {
    case protocol.Method_CHANNEL_ACK:
        return a.onChannelAck(commonPt.Payload)
    case protocol.Method_STREAM_CREATE:
        return a.onStreamCreate(commonPt.Payload)
    case protocol.Method_STREAM_ACK:
        return a.onStreamAck(commonPt.Payload)
    case protocol.Method_STREAM_CLOSE:
        fallthrough
    case protocol.Method_CHANNEL_CREATE:
        // 客户端不会接收到此消息，由Bytezero处理.
        fallthrough
    default:
        fmt.Printf("AppsChannel.handlePt - Method type %v no impliment\n", commonPt.Type)
        return protocol.ErrBZProtocol
    }
    // return nil
}

// channelCreate -
func (a *AppsChannel) channelCreate(sessionId string) (err error) {
    a.state = protocol.ChannelStateCreate
    a.sessionId = sessionId
    channelCreatePt := &protocol.ChannelCreatePt {
        AppId: []byte(a.app.AppId()),
        DeviceId: []byte(a.app.DeviceId()),
        SessionId: []byte(sessionId),
    }
    mByte, err := protocol.Marshal(channelCreatePt)
    if err != nil {
        return err
    }
    return a.Send(mByte)
}

// channelClose -
func (a *AppsChannel) channelClose() (err error) {
    return a.TCPConn.Close()
}

// onChannelAck - code, message, chanId.
func (a *AppsChannel) onChannelAck(payload []byte) error {
    channelAckPt := &protocol.ChannelAckPt{}
    if err := protocol.Unmarshal(payload, channelAckPt); err != nil {
        return err
    }
    if channelAckPt.Code == protocol.ErrCode_success {
        a.cid = channelAckPt.Id
        a.app.OnSuccess(a.cid)
        a.state = protocol.ChannelStateOpen
    } else {
        a.app.OnError(int(channelAckPt.Code), channelAckPt.Message)
        a.state = protocol.ChannelStateFailed
    }
    a.ack <- channelAckPt.Code
    return nil
}

// StreamCreate -
func (a *AppsChannel) StreamCreate(sid protocol.StreamId, observer StreamObserver) (StreamHandle, error) {
    a.l.Lock()
    defer a.l.Unlock()
    streamHandle, ok := a.m[sid]
    if ok {
        return streamHandle, nil
    }

    streamHandle = NewAppsStream(sid, a, observer)
    if err := streamHandle.Create(); err != nil {
        return nil, err
    }
    a.m[sid] = streamHandle
    a.w[sid] = make(chan protocol.ErrCode)
    // wait.
    if err := a.wait(a.w[sid], waitTimeoutForAck, "stream", sid); err != nil {
        return nil, err
    }
    return streamHandle, nil
}

// onStreamCreate -
func (a *AppsChannel) onStreamCreate(payload []byte) error {
    // 请求.
    streamCreatePt := &protocol.StreamCreatePt{}
    if err := protocol.Unmarshal(payload, streamCreatePt); err != nil {
        return err
    }
    // 通知.
    a.l.Lock()
    defer a.l.Unlock()
    streamHandle := NewAppsStream(streamCreatePt.Id, a, nil)
    var code protocol.ErrCode = protocol.ErrCode_success
    var err error = fmt.Errorf("ok")
    code, err = a.app.onStream(streamHandle);
    if code == protocol.ErrCode_success {
        a.m[streamCreatePt.Id] = streamHandle
    }
    // 响应.
    fmt.Printf("AppsChannel.onStreamCreate - create %d, code %d, message %s \n", streamHandle.Id, code, err.Error())
    return streamHandle.Ack(code, err.Error())
}

// onStreamAck - code, message, chanId.
func (a *AppsChannel) onStreamAck(payload []byte) error {
    streamAckPt := &protocol.StreamAckPt{}
    if err := protocol.Unmarshal(payload, streamAckPt); err != nil {
        return err
    }

    a.l.Lock()
    defer a.l.Unlock()
    if stream, ok := a.m[streamAckPt.Id]; ok {
        err := stream.OnAck(streamAckPt)
        if r, ok := a.w[streamAckPt.Id]; ok {
            r <- streamAckPt.Code
        }
        return err
    }
    return fmt.Errorf("Stream Id<%d> no exist, ack code %d, message %s", streamAckPt.Id, streamAckPt.Code, streamAckPt.Message)
}
