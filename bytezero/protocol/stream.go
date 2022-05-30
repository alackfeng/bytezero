package protocol

import (
	"encoding/binary"
	"fmt"
)

// StreamState -
type StreamState uint32
const (
    StreamStateNone StreamState = iota
    StreamStateCreate
    StreamStateOpen
    StreamStateFailed
    StreamStateClosing
    StreamStateClosed
    StreamStateMax
)

// StreamId -
type StreamId uint32

// String -
func (s StreamId) String() string {
    return fmt.Sprintf("#%d", s)
}


// Serializable -
type Serializable interface {
    Len() int
    Pack(buf []byte, i int) int
    UnPack(buf []byte, i int) int
}

// StreamVer -
type StreamVer uint8
const (
    StreamVerNone       StreamVer = 0x0
    StreamVerExtra      StreamVer = 0x1 << 1
    StreamVerReserved   StreamVer = 0x1 << 2
    StreamVeAll         StreamVer = StreamVerExtra | StreamVerReserved
)

var _ Serializable = (*StreamVer)(nil)

// Match -
func (s StreamVer) Match(v StreamVer) bool {
    return s & v == v
}

// Len -
func (c *StreamVer) Len() int {
    return 1
}

// Pack -
func (c *StreamVer) Pack(buf []byte, i int) int {
    buf[i] = byte(*c)
    return 1
}

// UnPack -
func (c *StreamVer) UnPack(buf []byte, i int) int {
    *c = StreamVer(buf[i])
    return 1
}

/////////////////////////StreamCreatePt++++++++++++++++++++++++++++++++
// StreamCreatePt -
type StreamCreatePt struct {
    Ver StreamVer `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"` // Stream Protocol Version.
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
    Extra []byte `form:"Extra" json:"Extra" xml:"Extra" bson:"Extra" binding:"required"` // stream extra info for use.
}
var _ BZProtocol = (*StreamCreatePt)(nil)


// NewStreamCreatePt -
func NewStreamCreatePt() *StreamCreatePt {
    return &StreamCreatePt{
        Ver: StreamVerNone,
    }
}

// Type -
func (c *StreamCreatePt) Type() Method {
    return Method_STREAM_CREATE
}

// Len -
func (c *StreamCreatePt) Len() int {
    l := 1 + 4 + 4 // default.
    if c.Ver.Match(StreamVerExtra) {
        l += 4 + len(c.Extra)
    }
    return l
}

// String -
func (c *StreamCreatePt) String() string {
    return fmt.Sprintf("Stream#%d", c.Id)
}

// Unmarshal -
func (c *StreamCreatePt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    var i uint32 = 0
    c.Ver = StreamVer(buf[i]); i += 1
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    if c.Ver.Match(StreamVerExtra) {
        l := binary.BigEndian.Uint32(buf[i:]); i += 4
        fmt.Println("StreamCreatePt::Unmarshal - extra, length", l)
        c.Extra = buf[i:i+l]
    }
    return nil
}

// Marshal -
func (c *StreamCreatePt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    i := 0
    buf[i] = byte(c.Ver); i += 1 // Stream Protocol Version.
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Od)); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Id)); i += 4
    if c.Ver.Match(StreamVerExtra) {
        fmt.Println("StreamCreatePt::Marshal - extra, length", len(c.Extra))
        binary.BigEndian.PutUint32(buf[i:], uint32(len(c.Extra))); i += 4
        ByteCopy(buf, i, c.Extra, 0);
    } else {
        fmt.Println("StreamCreatePt::Marshal - no extra.")
    }
    return buf, nil
}


/////////////////////////StreamAckPt++++++++++++++++++++++++++++++++
// StreamAckPt -
type StreamAckPt struct {
    Code ErrCode `form:"Code" json:"Code" xml:"Code" bson:"Code" binding:"required"` // ack Code.
    Message []byte `form:"Message" json:"Message" xml:"Message" bson:"Message" binding:"required"` // ack Message.
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
}
var _ BZProtocol = (*StreamAckPt)(nil)


// NewStreamAckPt -
func NewStreamAckPt() *StreamAckPt {
    return &StreamAckPt{}
}

// Type -
func (c *StreamAckPt) Type() Method {
    return Method_STREAM_ACK
}

// Len -
func (c *StreamAckPt) Len() int {
    return 4 + 4 + 4 + 4 + len(c.Message)
}

// String -
func (c *StreamAckPt) String() string {
    return fmt.Sprintf("Stream#%d Code.%d, Message.%s", c.Id, c.Code, c.Message)
}

// Unmarshal -
func (c *StreamAckPt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    var i uint32 = 0
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Code = ErrCode(binary.BigEndian.Uint32(buf[i:])); i += 4

    lc := binary.BigEndian.Uint32(buf[i:]); i += 4
    c.Message = buf[i:i+lc]; i += lc
    return nil
}

// Marshal -
func (c *StreamAckPt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    i := 0
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Od)); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Id)); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Code)); i += 4

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.Message))); i += 4
    ByteCopy(buf, i, c.Message, 0); i += len(c.Message)
    return buf, nil
}

/////////////////////////StreamClosePt++++++++++++++++++++++++++++++++
// StreamClosePt -
type StreamClosePt struct {
    Ver StreamVer `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"` // Stream Protocol Version.
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
    Extra []byte `form:"Extra" json:"Extra" xml:"Extra" bson:"Extra" binding:"required"` // stream extra info for use.
}
var _ BZProtocol = (*StreamClosePt)(nil)


// NewStreamClosePt -
func NewStreamClosePt() *StreamClosePt {
    return &StreamClosePt{}
}

// Type -
func (c *StreamClosePt) Type() Method {
    return Method_STREAM_CLOSE
}

// Len -
func (c *StreamClosePt) Len() int {
    l := 1 + 4 + 4
    if c.Ver.Match(StreamVerExtra) {
        l += 4 + len(c.Extra)
    }
    return l
}

// String -
func (c *StreamClosePt) String() string {
    return fmt.Sprintf("Channel#%dStream#%d Close", c.Od, c.Id)
}

// Unmarshal -
func (c *StreamClosePt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    var i uint32 = 0
    c.Ver = StreamVer(buf[i]); i += 1
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    if c.Ver.Match(StreamVerExtra) {
        l := binary.BigEndian.Uint32(buf[i:]); i += 4
        c.Extra = buf[i:i+l]
    }
    return nil
}

// Marshal -
func (c *StreamClosePt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    i := 0
    buf[i] = byte(c.Ver); i += 1
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Od)); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Id)); i += 4
    if c.Ver.Match(StreamVerExtra) {
        binary.BigEndian.PutUint32(buf[i:], uint32(len(c.Extra))); i += 4
        ByteCopy(buf, i, c.Extra, 0);
    }
    return buf, nil
}


/////////////////////////StreamDataPt++++++++++++++++++++++++++++++++
// StreamDataPt -
type StreamDataPt struct {
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
    Length uint32  `form:"Length" json:"Length" xml:"Length" bson:"Length" binding:"required"` // data length.
    Binary Boolean `form:"Binary" json:"Binary" xml:"Binary" bson:"Binary" binding:"required"` // binary or text.
    Data []byte `form:"Data" json:"Data" xml:"Data" bson:"Data" binding:"required"` // data buffer.
}
var _ BZProtocol = (*StreamDataPt)(nil)


// NewStreamDataPt -
func NewStreamDataPt() *StreamDataPt {
    return &StreamDataPt{}
}

// Type -
func (c *StreamDataPt) Type() Method {
    return Method_STREAM_DATA
}

// Len -
func (c *StreamDataPt) Len() int {
    return 4 + 4 + 1 + 4 + int(c.Length)
}

// String -
func (c *StreamDataPt) String() string {
    return fmt.Sprintf("Channel#%dSteam#%d Binary(%d) Length(%d)", c.Od, c.Id, c.Binary, c.Length)
}

// Unmarshal -
func (c *StreamDataPt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    var i uint32 = 0
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Length = binary.BigEndian.Uint32(buf[i:]); i += 4
    c.Binary = Boolean(buf[i]); i += 1
    c.Data = buf[i:i+c.Length]; // i += c.Length
    return nil
}

// Marshal -
func (c *StreamDataPt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    i := 0
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Od)); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Id)); i += 4
    binary.BigEndian.PutUint32(buf[i:], c.Length); i += 4
    buf[i] = byte(c.Binary); i += 1
    ByteCopy(buf, i, c.Data, 0); // i += int(c.Length)
    return buf, nil
}
