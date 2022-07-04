package protocol

import (
	"encoding/binary"
	"fmt"
)

// ChannelState -
type ChannelState uint32
const (
    ChannelStateNone ChannelState = iota
    ChannelStateConnecting
    ChannelStateCreate
    ChannelStateOpen
    ChannelStateFailed
    ChannelStateClosing
    ChannelStateClosed
    ChannelStateMax
)

// String -
func (c ChannelState) String() string {
    switch c {
    case ChannelStateConnecting: return "Connecting"
    case ChannelStateCreate: return "Create"
    case ChannelStateOpen: return "Open"
    case ChannelStateFailed: return "Failed"
    case ChannelStateClosing: return "Closing"
    case ChannelStateClosed: return "Closed"
    case ChannelStateMax: return "Max"
    }
    return "None"
}

/////////////////////////ChannelCreatePt++++++++++++++++++++++++++++++++
// ChannelCreatePt -
type ChannelCreatePt struct {
    Ver BridgeVer `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"` // Bridge Protocol Version Bits.
    OS OSType `form:"OS" json:"OS" xml:"OS" bson:"OS" binding:"required"`
    AppId []byte `form:"AppId" json:"AppId" xml:"AppId" bson:"AppId" binding:"required"`
    DeviceId []byte `form:"DeviceId" json:"DeviceId" xml:"DeviceId" bson:"DeviceId" binding:"required"`
    SessionId []byte `form:"SessionId" json:"SessionId" xml:"SessionId" bson:"SessionId" binding:"required"`
    User []byte `form:"User" json:"User" xml:"User" bson:"User" binding:"required"`
    Sign []byte `form:"Sign" json:"Sign" xml:"Sign" bson:"Sign" binding:"required"`
}
var _ BZProtocol = (*ChannelCreatePt)(nil)

// NewChannelCreatePb -
func NewChannelCreatePb(t OSType) *ChannelCreatePt {
    return &ChannelCreatePt{OS: t}
}

// FieldsSign -
func (c *ChannelCreatePt) FieldsSign() []byte {
    l := 1 + 1 + 4 + len(c.AppId) + 4 + len(c.DeviceId) + 4 + len(c.SessionId)
    buf := make([]byte, l)

    i := 0
    buf[i] = byte(c.Ver); i += 1
    buf[i] = byte(c.OS); i += 1

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.AppId))); i += 4
    ByteCopy(buf, i, c.AppId, 0); i += len(c.AppId)

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.DeviceId))); i += 4
    ByteCopy(buf, i, c.DeviceId, 0); i += len(c.DeviceId)

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.SessionId))); i += 4
    ByteCopy(buf, i, c.SessionId, 0); i += len(c.SessionId)
    return buf
}

// Type -
func (c *ChannelCreatePt) Type() Method {
    return Method_CHANNEL_CREATE
}

// Len -
func (c *ChannelCreatePt) Len() int {
    return 1 + 1 + 4 + len(c.AppId) + 4 + len(c.DeviceId) + 4 + len(c.SessionId) + 4 + len(c.User) + 4 + len(c.Sign)
}

// String -
func (c *ChannelCreatePt) String() string {
    return fmt.Sprintf("AppId<%s> DeviceId<%s> SessionId<%s> OS<%v> Sign<%s:%s>", c.AppId, c.DeviceId, c.SessionId, c.OS, c.User, c.Sign)
}

// Unmarshal -
func (c *ChannelCreatePt) Unmarshal(buf []byte) error {
    var i uint32 = 0
    c.Ver = BridgeVer(buf[i]); i += 1
    c.OS = OSType(buf[i]); i += 1

    la := binary.BigEndian.Uint32(buf[i:]); i += 4
    c.AppId = buf[i:i+la]; i += la

    lb := binary.BigEndian.Uint32(buf[i:]); i += 4
    c.DeviceId = buf[i:i+lb]; i += lb

    lc := binary.BigEndian.Uint32(buf[i:]); i += 4
    c.SessionId = buf[i:i+lc]; i += lc

    ld := binary.BigEndian.Uint32(buf[i:]); i += 4
    c.User = buf[i:i+ld]; i += ld

    le := binary.BigEndian.Uint32(buf[i:]); i += 4
    c.Sign = buf[i:i+le]; i += le
    return nil
}

// Marshal -
func (c *ChannelCreatePt) Marshal(buf []byte) ([]byte, error) {
    i := 0
    buf[i] = byte(c.Ver); i += 1
    buf[i] = byte(c.OS); i += 1

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.AppId))); i += 4
    ByteCopy(buf, i, c.AppId, 0); i += len(c.AppId)

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.DeviceId))); i += 4
    ByteCopy(buf, i, c.DeviceId, 0); i += len(c.DeviceId)

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.SessionId))); i += 4
    ByteCopy(buf, i, c.SessionId, 0); i += len(c.SessionId)

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.User))); i += 4
    ByteCopy(buf, i, c.User, 0); i += len(c.User)

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.Sign))); i += 4
    ByteCopy(buf, i, c.Sign, 0); i += len(c.Sign)
    return buf, nil
}


/////////////////////////ChannelAckPt++++++++++++++++++++++++++++++++
// ChannelAckPt -
type ChannelAckPt struct {
    Ver BridgeVer `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"` // Bridge Protocol Version Bits.
    Id ChannelId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // channel Id.
    Code ErrCode `form:"Code" json:"Code" xml:"Code" bson:"Code" binding:"required"` // ack Code.
    Message []byte `form:"Message" json:"Message" xml:"Message" bson:"Message" binding:"required"` // ack Message.
}
var _ BZProtocol = (*ChannelAckPt)(nil)

// NewChannelAckPt -
func NewChannelAckPt() *ChannelAckPt {
    return &ChannelAckPt{}
}

// Type -
func (c *ChannelAckPt) Type() Method {
    return Method_CHANNEL_ACK
}

// Len - Ver + Id + Code + Message
func (c *ChannelAckPt) Len() int {
    return 1 + 4 + 4 + 4 + len(c.Message)
}

// String -
func (c *ChannelAckPt) String() string {
    return fmt.Sprintf("Channel#%d Code.%d, Message.%s", c.Id, c.Code, c.Message)
}

// Unmarshal -
func (c *ChannelAckPt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    var i uint32 = 0
    c.Ver = BridgeVer(buf[i]); i += 1
    c.Id = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Code = ErrCode(binary.BigEndian.Uint32(buf[i:])); i += 4

    lc := binary.BigEndian.Uint32(buf[i:]); i += 4
    c.Message = buf[i:i+lc]; i += lc
    return nil
}

// Marshal -
func (c *ChannelAckPt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    i := 0
    buf[i] = byte(c.Ver); i += 1
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Id)); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Code)); i += 4

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.Message))); i += 4
    ByteCopy(buf, i, c.Message, 0); i += len(c.Message)
    return buf, nil
}




