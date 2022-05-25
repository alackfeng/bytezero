package protocol

import "errors"

var ErrNoEnoughtBufferLen = errors.New("No Enought buffer len")
var ErrNotBothOnline = errors.New("Channel connections not both online")
var ErrBufferNotAllWrite = errors.New("Buffer not all write to client")
var ErrBZProtocol = errors.New("BZProtocol Type no impliment")
var ErrBufferNotAllSent = errors.New("Buffer not sent all")


// ErrCode -
type ErrCode uint32

const (
    ErrCode_success             ErrCode = 0
    ErrCode_error               ErrCode = 1
    ErrCode_ack                 ErrCode = 100100
)
