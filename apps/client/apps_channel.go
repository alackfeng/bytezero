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
const waitTimeoutForAck = 30000 // ms.

// AppsChannel -
type AppsChannel struct {
    utils.BufferRead
    app Client
    *net.TCPConn
    address string
    sessionId string
    cid protocol.ChannelId
    ack chan protocol.ErrCode
    state protocol.ChannelState

    nextSid protocol.StreamId

    Observer ChannelObserver

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
        BufferRead: *utils.NewBufferRead(app.MaxRecvBufferLen()),
        ack: make(chan protocol.ErrCode),
        state: protocol.ChannelStateNone,
        m: make(map[protocol.StreamId]*AppsStream),
        w: make(map[protocol.StreamId]chan protocol.ErrCode),
    }
}

// RegisterObserver -
func (a *AppsChannel) RegisterObserver(observer ChannelObserver) {
    a.Observer = observer
}

// UnRegisterObserver -
func (a *AppsChannel) UnRegisterObserver() {
    a.Observer = nil
}

// Address -
func (a *AppsChannel) Address() string {
    return a.address
}

// Id -
func (a *AppsChannel) Id() protocol.ChannelId {
    return a.cid
}

// String -
func (a *AppsChannel) String() string {
    return fmt.Sprintf("<#%s#Channel#%d>", a.sessionId, a.cid)
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
    fmt.Printf("AppsChannel.Start - create connection channel sessionId<%s>, to <%s>.\n", sessionId, address)
    // 1. 建立Tcp连接.
    a.address = address
    if err := a.dial(); err != nil {
        return err
    }
    // 2. 启动消息接收.
    go a.handleRecevier()

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

// State -
func (a *AppsChannel) State() protocol.ChannelState {
    return a.state
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
    a.stateChanged(protocol.ChannelStateConnecting)
    a.TCPConn = tcpConn
    return nil
}

// stateChanged -
func (a *AppsChannel) stateChanged(state protocol.ChannelState) {
    a.state = state
    if a.Observer != nil {
        a.Observer.OnState(a.state)
    }
}

// wait -
func (a *AppsChannel) wait(ack chan protocol.ErrCode, ms int, label string, sid protocol.StreamId) error {
    duration := time.Duration(ms) * time.Millisecond
    ticker := time.NewTicker(duration)
    fmt.Printf("AppsChannel::wait - %s# timeout %v, <- chan %v\n", label, duration, ack)
    select {
    case code := <- ack:
        // fmt.Printf("AppsChannel::wait - %s# , ErrCode %d, Online %v\n", label, code, a.Online())
        if a.Online() {
            return nil
        }
        if label == "channel" {
            return fmt.Errorf("Channel#%s, ErrCode.%d, State.%v", a.sessionId, code, a.State())
        } else {
            return fmt.Errorf("Stream#%d, ErrCode.%d", sid, code)
        }
    case <- ticker.C:
        if label == "channel" {
            return fmt.Errorf("Channel#%s timeout, State.%v", a.sessionId, a.State())
        } else {
            return fmt.Errorf("Stream#%d timeout", sid)
        }
    }
}

// handleRecevier -
func (a *AppsChannel) handleRecevier() {
    if a.address == "" {
        fmt.Println("AppsChannel::handleRecevier - address is null.")
        return
    }

    // recvBufferLen := a.app.MaxRecvBufferLen()
    // buffer := make([]byte, recvBufferLen)
    // currTime := time.Now()
    for {
        // n, err := a.Read(buffer)
        // if err != nil {
        //     fmt.Printf("AppsChannel::handleRecevier error.%v.\n", err)
        //     break
        // }
        // if a.recvStat.Bytes == 0 {
        //     a.recvStat.Begin()
        //     fmt.Printf("AppsChannel::handleRecevier - begin. recv begin %v.\n", a.recvStat.Info())
        // }
        // a.recvStat.Inc(int64(n))
        // if n != recvBufferLen {
        //     fmt.Printf("AppsChannel::handleRecevier - recv buffer len %d not equal recvBufferLen, real %d.\n", recvBufferLen, n)
        // }
        // if time.Now().Sub(currTime).Milliseconds() > 1000 {
        //     currTime = time.Now()
        //     fmt.Printf("AppsChannel::handleRecevier - recv count %d, bps %s. send bps %s\n", a.recvStat.Count, utils.ByteSizeFormat(a.recvStat.Bps1s()), utils.ByteSizeFormat(a.sendStat.Bps1s()))
        // }

        if _, err := a.BufferRead.Read(func(b []byte) (int, error) {
            return a.TCPConn.Read(b)
        }); err != nil {
            logc.Errorf("AppsChannel::handleRecevier - read error.", err.Error())
            break
        }
        if a.BufferRead.Empty() {
            logc.Debugln("Connection handleRecevier - wait next.")
            continue
        }

        // 处理接收到消息.
        commonPt := &protocol.CommonPt{}
        if err := protocol.Unmarshal(a.BufferRead.Get(), commonPt); err != nil {
            fmt.Printf("AppsChannel::handleRecevicer - Unmarshal Buffer error.%v.\n", err.Error())
            if err == protocol.ErrNoFixedMe {
                logc.Errorln("AppsChannel::handleRecevier - Unmarshal error.", err.Error())
                break
            }
            a.BufferRead.Step()
            continue
        }
        a.BufferRead.Next(commonPt.Len())

        if err := a.handlePt(commonPt); err != nil {
            fmt.Printf("AppsChannel::handleRecevicer - handlePt error.%v.\n", err.Error())
            continue
        }

    }
    a.recvStat.End()
    fmt.Printf("AppsChannel::handleRecevicer - end... %v.\n", a.recvStat.InfoAll())
    a.stateChanged(protocol.ChannelStateClosed)
    a.ack <- protocol.ErrCode_ConnectionClosed

    a.app.ChannelClose(a.sessionId)
}

// handlePt -
func (a *AppsChannel) handlePt(commonPt *protocol.CommonPt) error {
    switch commonPt.Type {
    case protocol.Method_CHANNEL_ACK:
        return a.onChannelAck(commonPt)
    case protocol.Method_STREAM_CREATE:
        return a.onStreamCreate(commonPt)
    case protocol.Method_STREAM_ACK:
        return a.onStreamAck(commonPt)
    case protocol.Method_STREAM_DATA:
        return a.onStreamData(commonPt)
    case protocol.Method_STREAM_CLOSE:
        return a.onStreamClose(commonPt)
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
    a.stateChanged(protocol.ChannelStateCreate)
    fmt.Println("AppsChannel channelCreate - ", channelCreatePt, ", len ", len(mByte))
    return a.Send(mByte)
}

// channelClose -
func (a *AppsChannel) channelClose() (err error) {
    if a.state == protocol.ChannelStateClosed {
        return nil
    }
    return a.TCPConn.Close()
}

// onChannelAck - code, message, chanId.
func (a *AppsChannel) onChannelAck(commonPt *protocol.CommonPt) error {
    channelAckPt := &protocol.ChannelAckPt{}
    if err := commonPt.UnmarshalP(channelAckPt); err != nil {
        return err
    }
    if channelAckPt.Code == protocol.ErrCode_success {
        a.cid = channelAckPt.Id
        a.Observer.OnSuccess(a.cid)
        a.stateChanged(protocol.ChannelStateOpen)
    } else {
        a.Observer.OnError(int(channelAckPt.Code), string(channelAckPt.Message))
        a.stateChanged(protocol.ChannelStateFailed)
    }
    a.ack <- channelAckPt.Code
    return nil
}

// StreamCreate -
func (a *AppsChannel) StreamCreate(observer StreamObserver, extra []byte) (StreamHandle, error) {
    a.l.Lock()

    sid := a.streamNextId()
    streamHandle, ok := a.m[sid]
    if ok {
        a.l.Unlock()
        return streamHandle, nil
    }

    streamHandle = NewAppsStream(sid, a, observer, extra)
    if err := streamHandle.Create(); err != nil {
        a.l.Unlock()
        return nil, err
    }
    a.m[sid] = streamHandle
    a.w[sid] = make(chan protocol.ErrCode)
    a.l.Unlock()
    // wait.
    if err := a.wait(a.w[sid], waitTimeoutForAck, "stream", sid); err != nil {
        return nil, err
    }
    return streamHandle, nil
}

// StreamClose -
func (a *AppsChannel) StreamClose(sid protocol.StreamId, extra []byte) error {
    a.l.Lock()

    streamHandle, ok := a.m[sid]
    if !ok {
        a.l.Unlock()
        return nil
    }
    streamHandle.extra = extra
    if err := streamHandle.Close(); err != nil {
        a.l.Unlock()
        return err
    }
    a.w[sid] = make(chan protocol.ErrCode)
    a.l.Unlock()
    // wait.
    if err := a.wait(a.w[sid], waitTimeoutForAck, "stream", sid); err != nil {
        return err
    }
    a.streamRemove(sid)
    return nil
}

// streamNextId -
func (a *AppsChannel) streamNextId() protocol.StreamId {
    a.nextSid++
    return a.nextSid
}

// streamRemove -
func (a *AppsChannel) streamRemove(sid protocol.StreamId) error {
    a.l.Lock()
    delete(a.m, sid)
    delete(a.w, sid)
    a.l.Unlock()
    return nil
}

// onStreamCreate -
func (a *AppsChannel) onStreamCreate(commonPt *protocol.CommonPt) error {
    // 请求.
    streamCreatePt := &protocol.StreamCreatePt{}
    if err := commonPt.UnmarshalP(streamCreatePt); err != nil {
        logc.Errorf("AppsChannel.onStreamCreate - StreamCreatePt UnmarshalP error.%v.", err.Error())
        return err
    }
    // 通知.
    a.l.Lock()
    defer a.l.Unlock()
    streamHandle := NewAppsStream(streamCreatePt.Id, a, nil, streamCreatePt.Extra)
    var code protocol.ErrCode = protocol.ErrCode_streamCreateAck
    var err error = fmt.Errorf("ok")
    code, err = a.Observer.onStream(streamHandle);
    if code == protocol.ErrCode_success || code == protocol.ErrCode_streamCreateAck {
        a.m[streamCreatePt.Id] = streamHandle
        // update stream id for
        if a.nextSid < streamCreatePt.Id {
            a.nextSid = streamCreatePt.Id
        }
    }
    // 响应.
    fmt.Printf("AppsChannel.onStreamCreate - create %d, code %d, message %s \n", streamHandle.Id, code, err.Error())
    return streamHandle.Ack(code, err.Error())
}

// onStreamAck - code, message, chanId.
func (a *AppsChannel) onStreamAck(commonPt *protocol.CommonPt) error {
    streamAckPt := &protocol.StreamAckPt{}
    if err := commonPt.UnmarshalP(streamAckPt); err != nil {
        return err
    }
    fmt.Printf("AppsChannel.onStreamAck - StreamAckPt %v.\n", streamAckPt)

    a.l.Lock()
    defer a.l.Unlock()
    if stream, ok := a.m[streamAckPt.Id]; ok {
        err := stream.OnAck(streamAckPt)
        if r, ok := a.w[streamAckPt.Id]; ok {
            fmt.Printf("AppsChannel.onStreamAck - notify streamAckPt Code.%v <- chan %v.\n", streamAckPt.Code, r)
            r <- streamAckPt.Code
        }
        fmt.Printf("AppsChannel.onStreamAck - notify OnAck Code.%v.\n", streamAckPt.Code)
        return err
    }
    return fmt.Errorf("Stream Id<%d> no exist, ack code %d, message %s", streamAckPt.Id, streamAckPt.Code, streamAckPt.Message)
}

// onStreamData -
func (a *AppsChannel) onStreamData(commonPt *protocol.CommonPt) error {
    streamDataPt := &protocol.StreamDataPt{}
    if err := commonPt.UnmarshalP(streamDataPt); err != nil {
        return err
    }
    fmt.Printf("AppsChannel.onStreamData - StreamDataPt %v.\n", streamDataPt)

    a.l.Lock()
    defer a.l.Unlock()
    if stream, ok := a.m[streamDataPt.Id]; ok {
        stream.onData(streamDataPt)
        return nil
    }
    fmt.Printf("AppsChannel.onStreamData - %v not exist.\n", streamDataPt)
    return nil
}

// onStreamClose -
func (a *AppsChannel) onStreamClose(commonPt *protocol.CommonPt) error {
    streamClosePt := &protocol.StreamClosePt{}
    if err := commonPt.UnmarshalP(streamClosePt); err != nil {
        return err
    }
    fmt.Printf("AppsChannel.onStreamClose - StreamClosePt %v.\n", streamClosePt)

    a.l.Lock()
    defer a.l.Unlock()
    if stream, ok := a.m[streamClosePt.Id]; ok {
        err := stream.onClose(streamClosePt)

        a.streamRemove(streamClosePt.Id)
        return err
    }
    fmt.Printf("AppsChannel.onStreamClose - %v not exist.\n", streamClosePt)
    return nil
}

