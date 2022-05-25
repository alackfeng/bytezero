package protocol

import (
	"fmt"
)

// ChannelId -
type ChannelId uint32

// Channel Method -
type Method uint8

const (
    Method_NONE                Method = 0
    Method_CHANNEL_CREATE      Method = 1
    Method_CHANNEL_ACK         Method = 2
    Method_CHANNEL_CLOSE       Method = 3
    Method_STREAM_CREATE       Method = 4
    Method_STREAM_ACK          Method = 5
    Method_STREAM_DATA         Method = 6
    Method_STREAM_CLOSE        Method = 7
    Method_STREAM_REFRESH      Method = 8
    Method_STREAM_ERROR        Method = 9
)

var Method_name = map[int32]string {
    0: "NONE",
    1: "CHANNEL_CREATE",
    2: "CHANNEL_ACK",
    3: "CHANNEL_CLOSE",
    4: "STREAM_CREATE",
    5: "STREAM_ACK",
    6: "STREAM_DATA",
    7: "STREAM_CLOSE",
}

// Method String -
func (x Method) String() string {
    return Method_name[int32(x)]
}

// Unmarshal -
func Unmarshal(buf []byte, out interface{}) error {
    switch out.(type) {
    case *CommonPt:
        return out.(*CommonPt).Unmarshal(buf)
    default:
        return fmt.Errorf("Protocol Unmarshal %v not support ", out)
    }
}

// Marshal -
func Marshal(m interface{}) ([]byte, error) {
    if cc, ok := m.(BZProtocol); ok {
        return NewCommPb(cc.Type()).MarshalP(cc)
    }
    return nil, ErrBZProtocol
}
