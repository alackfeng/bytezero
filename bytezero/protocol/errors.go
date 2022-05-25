package protocol

import "errors"

var ErrNoEnoughtBufferLen = errors.New("No Enought buffer len")
var ErrNotBothOnline = errors.New("Channel connections not both online")
var ErrBufferNotAllWrite = errors.New("Buffer not all write to client")
var ErrBZProtocol = errors.New("BZProtocol Type no impliment")
