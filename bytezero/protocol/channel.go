package protocol

import (
	"encoding/binary"
)

/////////////////////////ChannelCreatePt++++++++++++++++++++++++++++++++
// ChannelCreatePt -
type ChannelCreatePt struct {
    AppId []byte `form:"AppId" json:"AppId" xml:"AppId" bson:"AppId" binding:"required"`
    DeviceId []byte `form:"DeviceId" json:"DeviceId" xml:"DeviceId" bson:"DeviceId" binding:"required"`
    SessionId []byte `form:"SessionId" json:"SessionId" xml:"SessionId" bson:"SessionId" binding:"required"`
    Sign []byte `form:"Sign" json:"Sign" xml:"Sign" bson:"Sign" binding:"required"`
}
var _ BZProtocol = (*ChannelCreatePt)(nil)

// NewChannelCreatePb -
func NewChannelCreatePb() *ChannelCreatePt {
    return &ChannelCreatePt{}
}

// Type -
func (c *ChannelCreatePt) Type() Method {
    return Method_CHANNEL_CREATE
}

// Len -
func (c *ChannelCreatePt) Len() int {
    return 4 + len(c.AppId) + 4 + len(c.DeviceId) + 4 + len(c.SessionId)
}

// Unmarshal -
func (c *ChannelCreatePt) Unmarshal(buf []byte) error {
    var i uint32 = 0
    la := binary.BigEndian.Uint32(buf[i:i+4]); i += 4
    lb := binary.BigEndian.Uint32(buf[i:i+4]); i += 4
    lc := binary.BigEndian.Uint32(buf[i:i+4]); i += 4
    c.AppId = buf[i:i+la]; i += la
    c.DeviceId = buf[i:i+lb]; i += lb
    c.SessionId = buf[i:i+lc]; i += lc
    return nil
}

// Marshal -
func (c *ChannelCreatePt) Marshal(buf []byte) ([]byte, error) {
    i := 0
    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.AppId))); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.DeviceId))); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.SessionId))); i += 4
    ByteCopy(buf, i, c.AppId, 0); i += len(c.AppId)
    ByteCopy(buf, i, c.DeviceId, 0); i += len(c.DeviceId)
    ByteCopy(buf, i, c.SessionId, 0); i += len(c.SessionId)
    return buf, nil
}


/////////////////////////ChannelAckPt++++++++++++++++++++++++++++++++
// ChannelAckPt -
type ChannelAckPt struct {
    Id ChannelId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // channel id.
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

// Len -
func (c *ChannelAckPt) Len() int {
    return 4
}

// Unmarshal -
func (c *ChannelAckPt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    c.Id = ChannelId(binary.BigEndian.Uint32(buf[0:4]))
    return nil
}

// Marshal -
func (c *ChannelAckPt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    binary.BigEndian.PutUint32(buf[:], uint32(c.Id))
    return buf, nil
}




