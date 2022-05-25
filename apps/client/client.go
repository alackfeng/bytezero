package client

import "github.com/alackfeng/bytezero/bytezero/protocol"

// BaseHandle
type BaseHandle interface {
    ConnectionChannel(sessionId string) (ChannelHandle, error)
    ChannelClose(sessionId string) (err error)
}

// Client -
type Client interface {
    BaseHandle
    ChannelObserver
    MaxRecvBufferLen() int
    AppId() string
    DeviceId() string
    SessionId() string
    TargetAddress() string
}

// ChannelSender -
type ChannelSender interface {
    Send([]byte) error
}

// ChannelObserver -
type ChannelObserver interface {
    OnSuccess(protocol.ChannelId)
    OnError(int, string)
    onStream(StreamHandle) (protocol.ErrCode, error)
}

// ChannelHandle -
type ChannelHandle interface {
    // channel operator.
    Start(address string, sessionId string) error
    Stop() error
    Online() bool

    // stream operator.
    StreamCreate(sid protocol.StreamId, observer StreamObserver) (StreamHandle, error)
}



// StreamObserver -
type StreamObserver interface {
    OnStreamSuccess(protocol.StreamId)
    OnStreamError(int, string)
}

// StreamHandle -
type StreamHandle interface {
    Create() error
    Close() error
    StreamId() protocol.StreamId
    SendData([]byte) error
    SendSignal([]byte) error
}
