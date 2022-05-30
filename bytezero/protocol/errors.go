package protocol

import "errors"

var ErrNoEnoughtBufferLen = errors.New("No Enought buffer len")
var ErrNotBothOnline = errors.New("Channel connections not both online")
var ErrBufferNotAllWrite = errors.New("Buffer not all write to client")
var ErrBZProtocol = errors.New("BZProtocol Type no impliment")
var ErrBufferNotAllSent = errors.New("Buffer not sent all")
var ErrPackBufferNotEnought = errors.New("Pack Buffer Not Enought")
var ErrNoFixedMe = errors.New("No Fixed Me")
var ErrNoLength = errors.New("No Length")
var ErrNoPayload = errors.New("No Payload")
var ErrNoSessionId = errors.New("No SessionId")
var ErrNoDeviceId = errors.New("No DeviceId")


// ErrCode -
type ErrCode uint32

const (
    ErrCode_success             ErrCode = 0
    ErrCode_error               ErrCode = 1

    ErrCode_ConnectionError     ErrCode = 100000
    ErrCode_ConnectionClosed    ErrCode = 100001
    ErrCode_ack                 ErrCode = 100100
    ErrCode_streamCreateAck     ErrCode = 100101
    ErrCode_streamCloseAck      ErrCode = 100102

)
