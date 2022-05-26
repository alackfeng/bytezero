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


/////////////////////////StreamCreatePt++++++++++++++++++++++++++++++++
// StreamCreatePt -
type StreamCreatePt struct {
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
}
var _ BZProtocol = (*StreamCreatePt)(nil)


// NewStreamCreatePt -
func NewStreamCreatePt() *StreamCreatePt {
    return &StreamCreatePt{}
}

// Type -
func (c *StreamCreatePt) Type() Method {
    return Method_STREAM_CREATE
}

// Len -
func (c *StreamCreatePt) Len() int {
    return 4 + 4
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
    i := 0
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    return nil
}

// Marshal -
func (c *StreamCreatePt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    i := 0
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Od)); i += 4
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Id)); i += 4
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



/////////////////////////StreamErrorPt++++++++++++++++++++++++++++++++
// StreamErrorPt -
type StreamErrorPt struct {
    Code ErrCode `form:"Code" json:"Code" xml:"Code" bson:"Code" binding:"required"` // ack Code.
    Message []byte `form:"Message" json:"Message" xml:"Message" bson:"Message" binding:"required"` // ack Message.
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
}
var _ BZProtocol = (*StreamErrorPt)(nil)


// NewStreamErrorPt -
func NewStreamErrorPt() *StreamErrorPt {
    return &StreamErrorPt{}
}

// Type -
func (c *StreamErrorPt) Type() Method {
    return Method_STREAM_ERROR
}

// Len -
func (c *StreamErrorPt) Len() int {
    return 4 + 4 + 4 + 4 + len(c.Message)
}

// String -
func (c *StreamErrorPt) String() string {
    return fmt.Sprintf("Stream#%d Error.%d, Message.%s", c.Id, c.Code, c.Message)
}

// Unmarshal -
func (c *StreamErrorPt) Unmarshal(buf []byte) error {
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
func (c *StreamErrorPt) Marshal(buf []byte) ([]byte, error) {
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
    Od ChannelId `form:"Od" json:"Od" xml:"Od" bson:"Od" binding:"required"` // Channel id.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
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
    return 4 + 4
}

// String -
func (c *StreamClosePt) String() string {
    return fmt.Sprintf("Stream#%d Close", c.Id)
}

// Unmarshal -
func (c *StreamClosePt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    c.Od = ChannelId(binary.BigEndian.Uint32(buf[0:4]))
    c.Id = StreamId(binary.BigEndian.Uint32(buf[0:4]))
    return nil
}

// Marshal -
func (c *StreamClosePt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    binary.BigEndian.PutUint32(buf[:], uint32(c.Od))
    binary.BigEndian.PutUint32(buf[:], uint32(c.Id))
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
    return 4 + 4 + 1 + 4 + 4 + len(c.Data)
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
    c.Data = buf[i:i+c.Length]; i += c.Length
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
    ByteCopy(buf, i, c.Data, 0); i += int(c.Length)
    return buf, nil
}
