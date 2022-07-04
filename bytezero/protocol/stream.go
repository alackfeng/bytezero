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

// String -
func (c StreamState) String() string {
    switch c {
    case StreamStateCreate: return "Create"
    case StreamStateOpen: return "Open"
    case StreamStateFailed: return "Failed"
    case StreamStateClosing: return "Closing"
    case StreamStateClosed: return "Closed"
    case StreamStateMax: return "Max"
    }
    return "None"
}

// StreamId -
type StreamId uint32

// String -
func (s StreamId) String() string {
    return fmt.Sprintf("Stream#%d", s)
}

/////////////////////////StreamCreatePt++++++++++++++++++++++++++++++++
// StreamCreatePt -
type StreamCreatePt struct {
    Ver BridgeVer `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"` // Bridge Protocol Version Bits.
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
    Timestamp uint64 `form:"Timestamp" json:"Timestamp" xml:"Timestamp" bson:"Timestamp" binding:"required"` // Timestamp.
    Extra []byte `form:"Extra" json:"Extra" xml:"Extra" bson:"Extra" binding:"required"` // stream extra info for use.
}
var _ BZProtocol = (*StreamCreatePt)(nil)


// NewStreamCreatePt -
func NewStreamCreatePt() *StreamCreatePt {
    return &StreamCreatePt{
        Ver: BridgeVerNone,
    }
}

// Type -
func (c *StreamCreatePt) Type() Method {
    return Method_STREAM_CREATE
}

// Len -
func (c *StreamCreatePt) Len() int {
    l := 1 + 4 + 4 + 8 // default.
    if c.Ver.Match(BridgeVerExtra) {
        l += 4 + len(c.Extra)
    }
    return l
}

// String -
func (c *StreamCreatePt) String() string {
    return fmt.Sprintf("Channel#%d+Stream#%d Timestamp.%d", c.Od, c.Id, c.Timestamp)
}

// Unmarshal -
func (c *StreamCreatePt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    var i uint32 = 0
    c.Ver = BridgeVer(buf[i]); i += 1
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Timestamp = binary.BigEndian.Uint64(buf[i:]); i += 8
    if c.Ver.Match(BridgeVerExtra) {
        l := binary.BigEndian.Uint32(buf[i:]); i += 4
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
    binary.BigEndian.PutUint64(buf[i:], c.Timestamp); i += 8
    if c.Ver.Match(BridgeVerExtra) {
        binary.BigEndian.PutUint32(buf[i:], uint32(len(c.Extra))); i += 4
        ByteCopy(buf, i, c.Extra, 0);
    }
    return buf, nil
}


/////////////////////////StreamAckPt++++++++++++++++++++++++++++++++
// StreamAckPt -
type StreamAckPt struct {
    Ver BridgeVer `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"` // Bridge Protocol Version Bits.
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
    Timestamp uint64 `form:"Timestamp" json:"Timestamp" xml:"Timestamp" bson:"Timestamp" binding:"required"` // Stream Create Timestamp.
    ArrivedTs uint64 `form:"ArrivedTs" json:"ArrivedTs" xml:"ArrivedTs" bson:"ArrivedTs" binding:"required"` // Stream Ack Response Arrived Timestamp.
    Code ErrCode `form:"Code" json:"Code" xml:"Code" bson:"Code" binding:"required"` // ack Code.
    Message []byte `form:"Message" json:"Message" xml:"Message" bson:"Message" binding:"required"` // ack Message.
    Extra []byte `form:"Extra" json:"Extra" xml:"Extra" bson:"Extra" binding:"required"` // ack Message.
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
    l := 1 + 4 + 4 + 8 + 8 + 4 + 4 + len(c.Message)
    if c.Ver.Match(BridgeVerExtra) {
        l += 4 + len(c.Extra)
    }
    return l
}

// String -
func (c *StreamAckPt) String() string {
    return fmt.Sprintf("Channel#%d+Stream#%d Code.%d, Message.%s", c.Od, c.Id, c.Code, c.Message)
}

// Unmarshal -
func (c *StreamAckPt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    var i uint32 = 0
    c.Ver = BridgeVer(buf[i]); i += 1
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Timestamp = binary.BigEndian.Uint64(buf[i:]); i += 8
    c.ArrivedTs = binary.BigEndian.Uint64(buf[i:]); i += 8

    c.Code = ErrCode(binary.BigEndian.Uint32(buf[i:])); i += 4

    lc := binary.BigEndian.Uint32(buf[i:]); i += 4
    c.Message = buf[i:i+lc]; i += lc

    if c.Ver.Match(BridgeVerExtra) {
        ld := binary.BigEndian.Uint32(buf[i:]); i += 4
        c.Extra = buf[i:i+ld]; i += ld
    }
    return nil
}

// Marshal -
func (c *StreamAckPt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    i := 0
    buf[i] = byte(c.Ver); i += 1
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Od)); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Id)); i += 4
    binary.BigEndian.PutUint64(buf[i:], c.Timestamp); i += 8
    binary.BigEndian.PutUint64(buf[i:], c.ArrivedTs); i += 8

    binary.BigEndian.PutUint32(buf[i:], uint32(c.Code)); i += 4

    binary.BigEndian.PutUint32(buf[i:], uint32(len(c.Message))); i += 4
    ByteCopy(buf, i, c.Message, 0); i += len(c.Message)

    if c.Ver.Match(BridgeVerExtra) {
        binary.BigEndian.PutUint32(buf[i:], uint32(len(c.Extra))); i += 4
        ByteCopy(buf, i, c.Extra, 0); i += len(c.Extra)
    }
    return buf, nil
}


/////////////////////////StreamClosePt++++++++++++++++++++++++++++++++
// StreamClosePt -
type StreamClosePt struct {
    Ver BridgeVer `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"` // Bridge Protocol Version Bits.
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
    if c.Ver.Match(BridgeVerExtra) {
        l += 4 + len(c.Extra)
    }
    return l
}

// String -
func (c *StreamClosePt) String() string {
    return fmt.Sprintf("Channel#%d+Stream#%d Close", c.Od, c.Id)
}

// Unmarshal -
func (c *StreamClosePt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    var i uint32 = 0
    c.Ver = BridgeVer(buf[i]); i += 1
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    if c.Ver.Match(BridgeVerExtra) {
        l := binary.BigEndian.Uint32(buf[i:]); i += 4
        c.Extra = buf[i:i+l]; i += l
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
    if c.Ver.Match(BridgeVerExtra) {
        binary.BigEndian.PutUint32(buf[i:], uint32(len(c.Extra))); i += 4
        ByteCopy(buf, i, c.Extra, 0); i += len(c.Extra)
    }
    return buf, nil
}


/////////////////////////StreamDataPt++++++++++++++++++++++++++++++++
// StreamDataPt -
type StreamDataPt struct {
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
    Timestamp uint64 `form:"Timestamp" json:"Timestamp" xml:"Timestamp" bson:"Timestamp" binding:"required"` // Timestamp.
    Total uint32 `form:"Total" json:"Total" xml:"Total" bson:"Total" binding:"required"` // data total.
    Offset uint32 `form:"Offset" json:"Offset" xml:"Offset" bson:"Offset" binding:"required"` // data offset.
    Length uint32 `form:"Length" json:"Length" xml:"Length" bson:"Length" binding:"required"` // data length.
    Ver BridgeVer `form:"Ver" json:"Ver" xml:"Ver" bson:"Ver" binding:"required"` // Bridge Protocol Version Bits.
    Binary Boolean `form:"Binary" json:"Binary" xml:"Binary" bson:"Binary" binding:"required"` // true: binary or false:text.
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
    return 4 + 4 + 8 + 4 + 4 + 4 + 1 + 1 + int(c.Length)
}

// String -
func (c *StreamDataPt) String() string {
    return fmt.Sprintf("Channel#%d+Steam#%d Binary(%d) Total(%d) Offset(%d) Timestamp(%d) Length(%d)", c.Od, c.Id, c.Binary, c.Total, c.Offset, c.Timestamp, c.Length)
}

// Unmarshal -
func (c *StreamDataPt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    var i uint32 = 0
    c.Ver = BridgeVer(buf[i]); i += 1
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Binary.From(buf[i]); i += 1
    c.Timestamp = binary.BigEndian.Uint64(buf[i:]); i += 8
    c.Total = binary.BigEndian.Uint32(buf[i:]); i += 4
    c.Offset = binary.BigEndian.Uint32(buf[i:]); i += 4
    c.Length = binary.BigEndian.Uint32(buf[i:]); i += 4
    if c.Data == nil {
        c.Data = make([]byte, c.Length)
    }
    ByteCopy(c.Data, 0, buf[i:i+c.Length], 0)
    return nil
}

// Marshal -
func (c *StreamDataPt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    i := 0
    buf[i] = byte(c.Ver); i += 1
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Od)); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Id)); i += 4
    buf[i] = byte(c.Binary); i += 1
    binary.BigEndian.PutUint64(buf[i:], c.Timestamp); i += 8
    binary.BigEndian.PutUint32(buf[i:], c.Total); i += 4
    binary.BigEndian.PutUint32(buf[i:], c.Offset); i += 4
    binary.BigEndian.PutUint32(buf[i:], c.Length); i += 4
    ByteCopy(buf, i, c.Data, 0); // i += int(c.Length)
    return buf, nil
}
