package protocol

import (
	"fmt"
)

// ChannelId -
type ChannelId uint32

// String -
func (s ChannelId) String() string {
    return fmt.Sprintf("Channel#%d", s)
}

var gloablChannelId ChannelId = 0
// GenChannelId -
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
    Method_MAX
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


// Serializable -
type Serializable interface {
    Len() int
    Pack(buf []byte, i int) int
    UnPack(buf []byte, i int) int
}

// BridgeVer -
type BridgeVer uint8
const (
    BridgeVerNone               BridgeVer = 0x0
    BridgeVerExtra              BridgeVer = 0x1 << 0
    BridgeVerReserved           BridgeVer = 0x1 << 7
    BridgeVerAll                BridgeVer = BridgeVerExtra | BridgeVerReserved
)

var _ Serializable = (*BridgeVer)(nil)

// Match -
func (s BridgeVer) Match(v BridgeVer) bool {
    return s & v == v
}

// Len -
func (c *BridgeVer) Len() int {
    return 1
}

// Pack -
func (c *BridgeVer) Pack(buf []byte, i int) int {
    buf[i] = byte(*c)
    return 1
}

// UnPack -
func (c *BridgeVer) UnPack(buf []byte, i int) int {
    *c = BridgeVer(buf[i])
    return 1
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
