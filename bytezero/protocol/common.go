package protocol

import (
	"encoding/binary"
	"fmt"
)

const FixedHead = 0xABCD

const CurrentVersion = VersionFirst

// BZProtocol -
type BZProtocol interface {
    Type() Method
    Len() int
    Unmarshal(buf []byte) error
    Marshal(buf []byte) ([]byte, error)
}

/////////////////////////HeadPt++++++++++++++++++++++++++++++++
// HeadPt.
type HeadPt struct {
    Fixed uint32 `form:"Fixed" json:"Fixed" xml:"Fixed" bson:"Fixed" binding:"required"`
    Ver VersionNumber `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"`
    Type Method `form:"Type" json:"Type" xml:"Type" bson:"Type" binding:"required"`
}

// NewHeadPb -
func NewHeadPb(method Method) *HeadPt {
    return &HeadPt{
        Fixed: uint32(FixedHead),
        Ver: CurrentVersion,
        Type: method,
    }
}

// // SetPayloadLength -
// func (c *HeadPt) SetPayloadLength(length int) *HeadPt {
//     c.Length = uint32(length)
//     return c
// }

// Len -
func (c *HeadPt) Len() int {
    // Fixed + Ver + Type
    return 4 + 2 + 1
}

// Pack -
func (c *HeadPt) Pack(buf []byte) error {
    binary.BigEndian.PutUint32(buf[0:4], c.Fixed)
    binary.BigEndian.PutUint16(buf[4:6], uint16(c.Ver))
    buf[6] = byte(c.Type)
    // binary.BigEndian.PutUint16(buf[7:11], uint16(c.Length))
    return nil
}

// UnPack -
func (c *HeadPt) UnPack(buf []byte) error {
    if len(buf) < c.Len() {
        return fmt.Errorf("HeadPt UnPack len %d less %d", len(buf), c.Len())
    }
    c.Fixed = binary.BigEndian.Uint32(buf[0:])
    if c.Fixed != FixedHead {
        return fmt.Errorf("No Fixed Head")
    }
    c.Ver = VersionNumber(binary.BigEndian.Uint16(buf[4:]))
    c.Type = Method(buf[6])
    return nil
}


/////////////////////////CommonPt++++++++++++++++++++++++++++++++
// CommonPt.
type CommonPt struct {
    HeadPt `form:"Head" json:"Head" xml:"Head" bson:"Head" binding:"required"`
    Length uint32 `form:"Length" json:"Length" xml:"Length" bson:"Length" binding:"required"`
    Payload []byte `form:"Payload" json:"Payload" xml:"Payload" bson:"Payload" binding:"required"`
}

// NewCommPb -
func NewCommPb(method Method) *CommonPt {
    return &CommonPt{
        HeadPt: HeadPt{
            Fixed: uint32(FixedHead),
            Ver: CurrentVersion,
            Type: method,
        },
    }
}

// Len -
func (c *CommonPt) Len() int {
    // Fixed + Ver + Type
    return c.HeadPt.Len() + 4 + int(c.Length)
}

// SetPayload -
func (c *CommonPt) SetPayload(b []byte) *CommonPt {
    c.Length = uint32(len(b))
    c.Payload = b
    return c
}

// Unmarshal -
func (c *CommonPt) Unmarshal(buf []byte) error {
    if err := c.HeadPt.UnPack(buf); err != nil {
        return err
    }
    l := len(buf)
    n := c.HeadPt.Len() - 1
    if n + 4 < l {
        return fmt.Errorf("No Length")
    }
    c.Length = binary.BigEndian.Uint32(buf[n:]); n += 4
    if n + int(c.Length) < l {
        return fmt.Errorf("No Payload")
    }
    c.Payload = buf[n:]
    return nil
}

// Marshal -
func (c *CommonPt) Marshal() ([]byte, error) {
    l := c.Len()
    buf := make([]byte, l)
    if err := c.HeadPt.Pack(buf); err != nil {
        return nil, nil
    }
    n := c.HeadPt.Len() - 1
    binary.BigEndian.PutUint32(buf[n:], c.Length); n += 4
    ByteCopy(buf, n, c.Payload, 0)
    return buf, nil
}

// MarshalP -
func (c *CommonPt) MarshalP(m BZProtocol)([]byte, error) {
    // Head.
    l := c.HeadPt.Len() + 4 + m.Len()
    buf := make([]byte, l)
    if err := c.HeadPt.Pack(buf[0:]); err != nil {
        return nil, err
    }
    n := c.HeadPt.Len() - 1
    // Length.
    binary.BigEndian.PutUint32(buf[n:], uint32(m.Len())); n += 4
    // Payload.
    if _, err := m.Marshal(buf[n:]); err != nil {
        return nil, err
    }
    return buf, nil
}


// ByteCopy -
func ByteCopy(dst []byte, dstOffset int, src []byte, srcOffset int) error {
    dstCount := len(dst) - dstOffset
    srcCount := len(src) - srcOffset
    if dstCount < srcCount {
        return fmt.Errorf("Memcpy dst len %d less src len %d", dstCount, srcCount)
    }
    for i := 0; i < srcCount; i++ {
        dst[i+dstOffset] = src[i+srcOffset]
    }
    return nil
}
