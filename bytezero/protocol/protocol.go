package protocol

import (
	"fmt"
)

// ChannelId -
type ChannelId uint32

var gloablChannelId ChannelId = 0
// Next -
func GenChannelId() ChannelId  {
    gloablChannelId++
    return gloablChannelId
}

// Method - channel.
type Method uint8

const (
    Method_NONE                Method = iota
    Method_CHANNEL_CREATE
    Method_CHANNEL_ACK
    Method_CHANNEL_CLOSE
    Method_STREAM_CREATE
    Method_STREAM_ACK
    Method_STREAM_DATA
    Method_STREAM_CLOSE
    Method_STREAM_REFRESH
    Method_STREAM_ERROR
)

// Method String -
func (x Method) String() string {
    switch x {
    case Method_CHANNEL_CREATE: return "CHANNEL_CREATE"
    case Method_CHANNEL_ACK:    return "CHANNEL_ACK"
    case Method_CHANNEL_CLOSE:  return "CHANNEL_CLOSE"
    case Method_STREAM_CREATE:  return "STREAM_CREATE"
    case Method_STREAM_ACK:     return "STREAM_ACK"
    case Method_STREAM_DATA:    return "STREAM_DATA"
    case Method_STREAM_CLOSE:   return "STREAM_CLOSE"
    case Method_STREAM_REFRESH: return "STREAM_REFRESH"
    case Method_STREAM_ERROR:   return "STREAM_ERROR"
    }
    return "NONE"
}

// Unmarshal -
func Unmarshal(buf []byte, out interface{}) error {
    switch out.(type) {
    case *CommonPt:
        return out.(*CommonPt).Unmarshal(buf)
    case BZProtocol:
        return out.(BZProtocol).Unmarshal(buf)
    default:
        return fmt.Errorf("Protocol Unmarshal %v not support ", out)
    }
}

// Marshal -
func Marshal(m interface{}) ([]byte, error) {
    switch m.(type) {
    case *CommonPt:
        return m.(*CommonPt).Marshal()
    case BZProtocol:
        cc, _ := m.(BZProtocol)
        return NewCommPb(cc.Type()).MarshalP(cc)
    }
    fmt.Println("Protocol.Marshal - type not found.", m)
    return nil, ErrBZProtocol
}
