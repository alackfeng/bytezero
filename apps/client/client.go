package client

import "github.com/alackfeng/bytezero/bytezero/protocol"

// BaseHandle
type BaseHandle interface {
    ConnectionChannel(sessionId string, observer ChannelObserver) (ChannelHandle, error)
    ChannelClose(sessionId string) (err error)
    ChannelCloseByHandle(ChannelHandle) (err error)
}

// Client -
type Client interface {
    BaseHandle
    MaxRecvBufferLen() int
    AppId() string
    DeviceId() string
    SessionId() string
    TargetAddress() string
}

// ChannelSender -
type ChannelSender interface {
    Send([]byte) error
    Id() protocol.ChannelId
}

// ChannelObserver -
type ChannelObserver interface {
    OnSuccess(protocol.ChannelId)
    OnError(int, string)
    OnState(protocol.ChannelState)
    onStream(StreamHandle) (protocol.ErrCode, error)
}

// ChannelHandle -
type ChannelHandle interface {
    // channel operator.
    Start(address string, sessionId string) error
    Stop() error
    Online() bool
    Id() protocol.ChannelId

    // stream operator.
    StreamCreate(observer StreamObserver, extra []byte) (StreamHandle, error)
    StreamClose(sid protocol.StreamId, extra []byte) error
}



// StreamObserver -
type StreamObserver interface {
    OnStreamSuccess(protocol.StreamId)
    OnStreamError(int, string)
    OnStreamData([]byte, protocol.Boolean)
    OnStreamClosing([]byte)
}

// StreamHandle -
type StreamHandle interface {
    Create() error
    Close() error
    RegisterObserver(StreamObserver)
    UnRegisterObserver()
    StreamId() protocol.StreamId
    ExtraInfo() []byte
    SendData([]byte) error
    SendSignal([]byte) error
}
