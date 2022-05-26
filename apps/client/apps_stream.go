package client

import (
	"fmt"

	"github.com/alackfeng/bytezero/bytezero/protocol"
)

// AppsStream -
type AppsStream struct {
    Id protocol.StreamId
    Sender ChannelSender
    Observer StreamObserver
    state protocol.StreamState
}
var _ StreamHandle = (*AppsStream)(nil)

// NewAppsStream -
func NewAppsStream(sid protocol.StreamId, sender ChannelSender , observer StreamObserver) *AppsStream {
    return &AppsStream{
        Id: sid,
        Sender: sender,
        Observer: observer,
    }
}

// Create -
func (a *AppsStream) Create() error {
    a.state = protocol.StreamStateCreate
    streamCreatePt := &protocol.StreamCreatePt{
        Od: a.Sender.Id(),
        Id: a.Id,
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
    return nil
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
    if streamAckPt.Code == protocol.ErrCode_success {
        a.state = protocol.StreamStateOpen
        a.Observer.OnStreamSuccess(a.Id)
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

// SendData -
func (a *AppsStream) SendData(m []byte) error {
    streamDataPt := &protocol.StreamDataPt{
        Od: a.Sender.Id(),
        Id: a.Id,
        Binary: protocol.BooleanTrue,
        Length: uint32(len(m)),
        Data: m,
    }
    mByte, err := protocol.Marshal(streamDataPt);
    if err != nil {
        return err
    }
    return a.Sender.Send(mByte)
}

// SendSignal -
func (a *AppsStream) SendSignal(m []byte) error {
    streamDataPt := &protocol.StreamDataPt{
        Od: a.Sender.Id(),
        Id: a.Id,
        Binary: protocol.BooleanFalse,
        Length: uint32(len(m)),
        Data: m,
    }
    mByte, err := protocol.Marshal(streamDataPt);
    if err != nil {
        return err
    }
    return a.Sender.Send(mByte)
}




