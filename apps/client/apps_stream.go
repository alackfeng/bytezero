package client

import (
	"fmt"
	"time"

	"github.com/alackfeng/bytezero/bytezero/protocol"
)

// AppsStream -
type AppsStream struct {
    Id protocol.StreamId
    Sender ChannelSender
    Observer StreamObserver
    state protocol.StreamState
    extra []byte // extra info.
}
var _ StreamHandle = (*AppsStream)(nil)

// NewAppsStream -
func NewAppsStream(sid protocol.StreamId, sender ChannelSender , observer StreamObserver, extra []byte) *AppsStream {
    return &AppsStream{
        Id: sid,
        Sender: sender,
        Observer: observer,
        extra: extra,
    }
}

// Create -
func (a *AppsStream) Create() error {
    a.state = protocol.StreamStateCreate
    streamCreatePt := &protocol.StreamCreatePt{
        Od: a.Sender.Id(),
        Id: a.Id,
    }
    if len(a.extra) != 0 {
        streamCreatePt.Ver = protocol.StreamVerExtra
        streamCreatePt.Extra = a.extra
    }
    mByte, err := protocol.Marshal(streamCreatePt)
    if err != nil {
        logc.Errorf("AppsStream.Create - StreamCreatePt Marshal error.%v", err.Error())
        return err
    }
    logc.Debugf("AppsStream.Create - Send to %v, buffer %v.", a.Sender, streamCreatePt)
    if a.Sender != nil {
        return a.Sender.Send(mByte)
    }
    logc.Errorf("AppsStream.Create - Stream Sender is null.")
    return fmt.Errorf("Stream Sender is null")
}

// Close -
func (a *AppsStream) Close() error {
    a.state = protocol.StreamStateClosing
    streamClosePt := &protocol.StreamClosePt{
        Od: a.Sender.Id(),
        Id: a.Id,
    }
    if len(a.extra) != 0 {
        streamClosePt.Ver = protocol.StreamVerExtra
        streamClosePt.Extra = a.extra
    }
    mByte, err := protocol.Marshal(streamClosePt)
    if err != nil {
        logc.Errorf("AppsStream.Close - StreamClosePt Marshal error.%v", err.Error())
        return err
    }
    logc.Debugf("AppsStream.Close - Send to %v, buffer %v.", a.Sender, streamClosePt)
    if a.Sender != nil {
        return a.Sender.Send(mByte)
    }
    logc.Errorf("AppsStream.Close - Stream Sender is null.")
    return fmt.Errorf("Stream Sender is null")
}

// RegisterObserver -
func (a *AppsStream) RegisterObserver(observer StreamObserver) {
    a.Observer = observer
}

// UnRegisterObserver -
func (a *AppsStream) UnRegisterObserver() {
    a.Observer = nil
}


// StreamId -
func (a *AppsStream) StreamId() protocol.StreamId {
    return a.Id
}

// ExtraInfo -
func (a *AppsStream)  ExtraInfo() []byte {
    return a.extra
}

// Ack -
func (a *AppsStream) Ack(code protocol.ErrCode, message string) error {
    streamAckPt := &protocol.StreamAckPt{
        Od: a.Sender.Id(),
        Id: a.Id,
        Code: protocol.ErrCode(code),
        Message: []byte(message),
    }
    mByte, err := protocol.Marshal(streamAckPt)
    if err != nil {
        logc.Errorf("AppsStream.Create - StreamAckPt Marshal error.%v", err.Error())
        return err
    }
    return a.Sender.Send(mByte)
}

// OnAck -
func (a *AppsStream) OnAck(streamAckPt *protocol.StreamAckPt) error {
    if streamAckPt.Code == protocol.ErrCode_streamCreateAck {
        a.state = protocol.StreamStateOpen
        a.Observer.OnStreamSuccess(a.Id)
    } else if streamAckPt.Code == protocol.ErrCode_streamCloseAck {
        a.state = protocol.StreamStateClosed
    } else {
        a.state = protocol.StreamStateFailed
        a.Observer.OnStreamError(int(streamAckPt.Code),  string(streamAckPt.Message))
    }
    return nil
}

// onData -
func (a *AppsStream) onData(streamDataPt *protocol.StreamDataPt) error {
    if a.Observer != nil {
        a.Observer.OnStreamData(streamDataPt.Data, streamDataPt.Binary)
    }
    return nil
}

// onClose -
func (a *AppsStream) onClose(streamClosePt *protocol.StreamClosePt) error {
    if a.Observer != nil {
        a.Observer.OnStreamClosing(streamClosePt.Extra)
    }
    streamAckPt := &protocol.StreamAckPt{
        Od: a.Sender.Id(),
        Id: a.Id,
        Code: protocol.ErrCode_streamCloseAck,
        Message: []byte("normal closed"),
    }
    mByte, err := protocol.Marshal(streamAckPt)
    if err != nil {
        logc.Errorf("AppsStream.Create - StreamAckPt Marshal error.%v", err.Error())
        return err
    }
    return a.Sender.Send(mByte)
}

// SendData -
func (a *AppsStream) SendData(m []byte) error {
    l := len(m)
    for i:=0; i<l; i += 1024 {
        n := i + 1024
        if n > l {
            n = l
        }
        data := m[i:n]
        streamDataPt := &protocol.StreamDataPt{
            Od: a.Sender.Id(),
            Id: a.Id,
            Binary: protocol.BooleanTrue,
            Timestamp: uint64(time.Now().UnixNano()),
            Total: uint32(l),
            Offset: uint32(i),
            Length: uint32(len(data)),
            Data: data,
        }
        mByte, err := protocol.Marshal(streamDataPt);
        if err != nil {
            return err
        }
        if err := a.Sender.Send(mByte); err != nil {
            return err
        }
    }
    return nil
}

// SendSignal -
func (a *AppsStream) SendSignal(m []byte) error {
    l := len(m)
    for i:=0; i<l; i += 1024 {
        n := i + 1024
        if n > l {
            n = l
        }
        data := m[i:n]
        streamDataPt := &protocol.StreamDataPt{
            Od: a.Sender.Id(),
            Id: a.Id,
            Binary: protocol.BooleanFalse,
            Timestamp: uint64(time.Now().UnixNano()),
            Total: uint32(l),
            Offset: uint32(i),
            Length: uint32(len(data)),
            Data: data,
        }
        mByte, err := protocol.Marshal(streamDataPt);
        if err != nil {
            return err
        }
        if err := a.Sender.Send(mByte); err != nil {
            return err
        }
    }
    return nil
}



