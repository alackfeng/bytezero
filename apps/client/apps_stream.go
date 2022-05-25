package client

import "github.com/alackfeng/bytezero/bytezero/protocol"

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
    streamCreatePt := protocol.StreamCreatePt{
        Id: a.Id,
    }
    mByte, err := protocol.Marshal(streamCreatePt)
    if err != nil {
        return err
    }
    if a.Sender != nil {
        return a.Sender.Send(mByte)
    }
    return nil
}

// Close -
func (a *AppsStream) Close() error {
    return nil
}


// StreamId -
func (a *AppsStream) StreamId() protocol.StreamId {
    return a.Id
}

// Ack -
func (a *AppsStream) Ack(code protocol.ErrCode, message string) error {
    streamAckPt := &protocol.StreamAckPt{
        Id: a.Id,
        Code: protocol.ErrCode(code),
        Message: message,
    }
    mByte, err := protocol.Marshal(streamAckPt)
    if err != nil {
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
        a.Observer.OnStreamError(int(streamAckPt.Code), streamAckPt.Message)
    }
    return nil
}

// SendData -
func (a *AppsStream) SendData(m []byte) error {
    streamDataPt := protocol.StreamDataPt{
        Id: a.Id,
        Binary: true,
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
    streamDataPt := protocol.StreamDataPt{
        Id: a.Id,
        Binary: false,
        Length: uint32(len(m)),
        Data: m,
    }
    mByte, err := protocol.Marshal(streamDataPt);
    if err != nil {
        return err
    }
    return a.Sender.Send(mByte)
}




