package protocol

import (
	"encoding/binary"
	"fmt"
)

const FixedHead = 0xABCD

const CurrentVersion = Version20220526

type Boolean uint8
const (
    BooleanFalse Boolean = iota
    BooleanTrue
)

// To -
func (b Boolean) To() bool {
    return b == '1'
}

// From -
func (b Boolean) From(bl bool) Boolean {
    if bl == true {
        b = '1'
    } else {
        b = '0'
    }
    return b
}

// String -
func (b Boolean) String() string {
    if b == BooleanTrue {
        return "true"
    }
    return "false"
}

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
    Fixed uint16 `form:"Fixed" json:"Fixed" xml:"Fixed" bson:"Fixed" binding:"required"`
    Ver VersionNumber `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"`
    Type Method `form:"Type" json:"Type" xml:"Type" bson:"Type" binding:"required"`
}

// NewHeadPb -
func NewHeadPb(method Method) *HeadPt {
    return &HeadPt{
        Fixed: uint16(FixedHead),
        Ver: CurrentVersion,
        Type: method,
    }
}

// Len -
func (c *HeadPt) Len() int {
    // Fixed + Ver + Type
    return 2 + 1 + 1
}

// Pack -
func (c *HeadPt) Pack(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrPackBufferNotEnought
    }
    i := 0
    binary.BigEndian.PutUint16(buf[i:], c.Fixed); i += 2
    buf[i] = byte(c.Ver); i += 1
    buf[i] = byte(c.Type); i += 1
    fmt.Printf("Pack: Fixed.0x%X, Ver.%v, Type.%v.\n", c.Fixed, c.Ver, c.Type)
    return nil
}

// UnPack -
func (c *HeadPt) UnPack(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrPackBufferNotEnought
    }
    i := 0
    c.Fixed = binary.BigEndian.Uint16(buf[i:]); i += 2
    if c.Fixed != FixedHead {
        return ErrNoFixedMe
    }
    c.Ver = VersionNumber(buf[i]); i += 1
    c.Type = Method(buf[i]); i += 1
    fmt.Printf("UnPack: Fixed.0x%X, Ver.%v, Type.%v.\n", c.Fixed, c.Ver, c.Type)
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
        HeadPt: *NewHeadPb(method),
    }
}

// Len -
func (c *CommonPt) Len() int {
    // Head + Length + Payload
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
    i := c.HeadPt.Len()
    if i + 4 > l {
        return ErrNoLength
    }
    c.Length = binary.BigEndian.Uint32(buf[i:]); i += 4
    if i + int(c.Length) > l {
        return ErrNoPayload
    }
    c.Payload = buf[i:i+int(c.Length)]
    return nil
}

// Marshal -
func (c *CommonPt) Marshal() ([]byte, error) {
    l := c.Len()
    buf := make([]byte, l)
    if err := c.HeadPt.Pack(buf); err != nil {
        return nil, nil
    }
    i := c.HeadPt.Len()
    binary.BigEndian.PutUint32(buf[i:], c.Length); i += 4
    ByteCopy(buf, i, c.Payload, 0)
    return buf, nil
}

// UnmarshalP - unmarshal payload to Type struct.
func (c *CommonPt) UnmarshalP(m interface{}) error {
    if pl, ok := m.(BZProtocol); ok {
        if pl.Type() != c.Type {
            return ErrBZProtocol
        }
        return pl.Unmarshal(c.Payload)
    }
    return ErrBZProtocol
}

// MarshalP -
func (c *CommonPt) MarshalP(m BZProtocol)([]byte, error) {
    // Head.
    l := c.HeadPt.Len() + 4 + m.Len()
    buf := make([]byte, l)
    if err := c.HeadPt.Pack(buf[0:]); err != nil {
        return nil, err
    }
    i := c.HeadPt.Len()
    // Length.
    binary.BigEndian.PutUint32(buf[i:], uint32(m.Len())); i += 4
    // Payload.
    if _, err := m.Marshal(buf[i:]); err != nil {
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
