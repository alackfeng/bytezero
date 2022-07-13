package protocol

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"time"
)

// CurrentMs -
func CurrentMs() uint64 {
    return uint64(time.Now().UnixNano() / 1e6)
}

// OSType -
type OSType uint8
const (
    OSTypeNone OSType = iota
    OSTypeLinux
    OSTypeMacOS
    OSTypeWindows
    OSTypeAndroid
    OSTypeIOS
    OSTypeWeb
    OSTypeMax
    // OSTypeMobile = OSTypeIOS | OSTypeAndroid
    // OSTypePC = OSTypeLinux | OSTypeMacOS | OSTypeWindows
)

// GOOSType -
func GOOSType() OSType {
    fmt.Println("------OSType ", runtime.GOOS)
    switch (runtime.GOOS) {
    case "linux": return OSTypeLinux
    case "darwin": return OSTypeMacOS
    case "windows": return OSTypeWindows
    case "android": return OSTypeAndroid
    case "ios": return OSTypeIOS
    case "web": return OSTypeWeb
    }
    return OSTypeNone
}

// String -
func (s OSType) String() string {
    switch (s) {
    case OSTypeLinux: return "Linux"
    case OSTypeMacOS: return "MacOS"
    case OSTypeWindows: return "Windows"
    case OSTypeAndroid: return "Android"
    case OSTypeIOS: return "IOS"
    case OSTypeWeb: return "Web"
    }
    return "None"
}

// Match -
func (s OSType) Match(v OSType) bool {
    return s & v == v
}

// Mobile -
func (s OSType) Mobile() bool {
    return s == OSTypeAndroid || s == OSTypeIOS
}

// PC -
func (s OSType) Desktop() bool {
    return s == OSTypeLinux || s == OSTypeMacOS || s == OSTypeWindows
}

const FixedHead uint16 = 0xABCD
const CurrentVersion = Version20220526

type Boolean uint8
const (
    BooleanFalse Boolean = iota
    BooleanTrue
)

// To -
func (b Boolean) To() byte {
    if b == BooleanFalse {
        return '0'
    }
    return '1'
}

// From -
func (b* Boolean) From(bl byte) {
    if bl == 0 {
        *b = BooleanFalse
    } else {
        *b = BooleanTrue
    }
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
    Timestamp uint64 `form:"Timestamp" json:"Timestamp" xml:"Timestamp" bson:"Timestamp" binding:"required"`
}

// NewHeadPb -
func NewHeadPb(method Method) *HeadPt {
    return &HeadPt{
        Fixed: FixedHead,
        Ver: CurrentVersion,
        Type: method,
        Timestamp: CurrentMs(),
    }
}

// Len - 12=Fixed + Ver + Type + Timestamp.
func (c *HeadPt) Len() int {
    return 2 + 1 + 1 + 8
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
    binary.BigEndian.PutUint64(buf[i:], c.Timestamp); i += 8
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
    if c.Type <= Method_NONE || c.Type >= Method_MAX {
        return ErrNoMethodType
    }
    c.Timestamp = binary.BigEndian.Uint64(buf[i:]); i += 8
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

// Ts -
func (c *CommonPt) Ts() uint64 {
    return c.HeadPt.Timestamp
}

// Len - Head + Length + Payload
func (c *CommonPt) Len() int {
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
    if c.Payload == nil {
        c.Payload = make([]byte, c.Length)
    }
    ByteCopy(c.Payload, 0, buf[i:i+int(c.Length)], 0)
    return nil
}

// Marshal -
func (c *CommonPt) Marshal() ([]byte, error) {
    l := c.Len()
    buf := make([]byte, l)
    if err := c.HeadPt.Pack(buf); err != nil {
        return nil, err
    }
    if c.Length != uint32(len(c.Payload)) { // Length is Payload len.
        return nil, ErrNoPayloadLen
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
            fmt.Printf("CommonPt.UnmarshalP type %v not equal %v\n", pl.Type(), c.Type)
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
