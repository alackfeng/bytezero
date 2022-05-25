package protocol

import (
	"encoding/binary"
	"fmt"
)

// StreamState -
type StreamState uint32
const (
    StreamStateNone StreamState = 0x0
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
    return 4
}

// Unmarshal -
func (c *StreamCreatePt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    c.Id = StreamId(binary.BigEndian.Uint32(buf[0:4]))
    return nil
}

// Marshal -
func (c *StreamCreatePt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    binary.BigEndian.PutUint32(buf[:], uint32(c.Id))
    return buf, nil
}


/////////////////////////StreamAckPt++++++++++++++++++++++++++++++++
// StreamAckPt -
type StreamAckPt struct {
    Code ErrCode `form:"Code" json:"Code" xml:"Code" bson:"Code" binding:"required"` // ack Code.
    Message string `form:"Message" json:"Message" xml:"Message" bson:"Message" binding:"required"` // ack Message.
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
}
var _ BZProtocol = (*StreamAckPt)(nil)


// NewStreamAckPt -
func NewStreamAckPt() *StreamAckPt {
    return &StreamAckPt{}
}

// Type -
func (c *StreamAckPt) Type() Method {
    return Method_STREAM_CREATE
}

// Len -
func (c *StreamAckPt) Len() int {
    return 4
}

// Unmarshal -
func (c *StreamAckPt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    c.Id = StreamId(binary.BigEndian.Uint32(buf[0:4]))
    return nil
}

// Marshal -
func (c *StreamAckPt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    binary.BigEndian.PutUint32(buf[:], uint32(c.Id))
    return buf, nil
}



/////////////////////////StreamErrorPt++++++++++++++++++++++++++++++++
// StreamErrorPt -
type StreamErrorPt struct {
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
    return 4
}

// Unmarshal -
func (c *StreamErrorPt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    c.Id = StreamId(binary.BigEndian.Uint32(buf[0:4]))
    return nil
}

// Marshal -
func (c *StreamErrorPt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    binary.BigEndian.PutUint32(buf[:], uint32(c.Id))
    return buf, nil
}


/////////////////////////StreamClosePt++++++++++++++++++++++++++++++++
// StreamClosePt -
type StreamClosePt struct {
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
}
var _ BZProtocol = (*StreamClosePt)(nil)


// NewStreamClosePt -
func NewStreamClosePt() *StreamClosePt {
    return &StreamClosePt{}
}

// Type -
func (c *StreamClosePt) Type() Method {
    return Method_STREAM_ERROR
}

// Len -
func (c *StreamClosePt) Len() int {
    return 4
}

// Unmarshal -
func (c *StreamClosePt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    c.Id = StreamId(binary.BigEndian.Uint32(buf[0:4]))
    return nil
}

// Marshal -
func (c *StreamClosePt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    binary.BigEndian.PutUint32(buf[:], uint32(c.Id))
    return buf, nil
}


/////////////////////////StreamDataPt++++++++++++++++++++++++++++++++
// StreamDataPt -
type StreamDataPt struct {
    Id StreamId `form:"Id" json:"Id" xml:"Id" bson:"Id" binding:"required"` // stream id.
    Length uint32  `form:"Length" json:"Length" xml:"Length" bson:"Length" binding:"required"` // data length.
    Binary bool `form:"Binary" json:"Binary" xml:"Binary" bson:"Binary" binding:"required"` // binary or text.
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
    return 4 + 1 + 4 + len(c.Data)
}

// Unmarshal -
func (c *StreamDataPt) Unmarshal(buf []byte) error {
    if len(buf) < c.Len() {
        return ErrNoEnoughtBufferLen
    }
    i := 0
    c.Id = StreamId(binary.BigEndian.Uint32(buf[i:])); i += 4
    c.Binary = bool(buf[i] == '1'); i += 1
    c.Length = binary.BigEndian.Uint32(buf[i:]); i += 4
    if c.Data == nil {
        c.Data = make([]byte, c.Length)
    }
    ByteCopy(c.Data, 0, buf[i:i+int(c.Length)], int(c.Length))
    return nil
}

// Marshal -
func (c *StreamDataPt) Marshal(buf []byte) ([]byte, error) {
    if len(buf) < c.Len() {
        return buf, ErrNoEnoughtBufferLen
    }
    i := 0
    binary.BigEndian.PutUint32(buf[i:], uint32(c.Id)); i += 4
    if c.Binary {
        buf[i] = '1'
    } else {
        buf[i] = '0'
    }
    i += 1
    binary.BigEndian.PutUint32(buf[i:], c.Length); i += 4
    ByteCopy(buf, i, c.Data, 0)
    return buf, nil
}
