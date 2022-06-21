package cores

import (
	"fmt"
	"net"
	"time"

	bz "github.com/alackfeng/bytezero/bytezero"
	"github.com/alackfeng/bytezero/bytezero/protocol"
	"github.com/alackfeng/bytezero/cores/utils"
)

const defaultMaxBufferLen = 1538

// Connection -
type Connection struct {
    utils.BufferRead
    *net.TCPConn
    bzn bz.BZNet
    maxBufferLen int
    quit bool
    durationMs int64

    // Info.
    DeviceId string
    SessionId string

    data chan []byte
}

var _ bz.BZNetReceiver = (*Connection)(nil)

// NewConnection -
func NewConnection(bzn bz.BZNet, c *net.TCPConn) *Connection {
    return &Connection{
        TCPConn: c,
        bzn: bzn,
        maxBufferLen: defaultMaxBufferLen*10,
        BufferRead: *utils.NewBufferRead(defaultMaxBufferLen*10),
        data: make(chan []byte, 10240),
    }
}

// Main -
func (c *Connection) Main()*Connection {
    go c.handleTrans()
    go c.handleRecevier()
    return c
}

// ChannId -
func (c *Connection) ChannId() string {
    return c.SessionId
}

// Id -
func (c *Connection) Id() string {
    return c.RemoteAddr().String()
}

// Set -
func (c *Connection) Set(info *protocol.ChannelCreatePt) *Connection {
    c.DeviceId = string(info.DeviceId)
    c.SessionId = string(info.SessionId)
    return c
}

// Check -
func (c *Connection) Check() error {
    if c.SessionId == "" {
        return protocol.ErrNoSessionId
    } else if c.DeviceId == "" {
        return protocol.ErrNoDeviceId
    }
    return nil
}

// Equal -
func (c *Connection) Equals(o *Connection) bool {
    return c.DeviceId == o.DeviceId
}

// Quit -
func (c *Connection) Quit() {
    c.Close()
    c.quit = true
}

func (c Connection) String() string {
    return fmt.Sprintf("Connection#[SessionId<%s>, DeviceId<%s>]", c.SessionId, c.DeviceId)
}


// Transit - to connection.
func (c *Connection) Transit(buf []byte) error {
    // return c.Send(buf)
    c.data <- buf
    return nil
}

// Send -
func (c *Connection) Send(buf []byte) error {
    n, err := c.Write(buf)
    if err != nil {
        return err
    }
    if n != len(buf) {
        return protocol.ErrBufferNotAllWrite
    }
    return nil
}

// handleTrans -
func (c *Connection) handleTrans() {
    for {
        select {
        case d, ok := <- c.data:
            if !ok {
                return
            }
            c.Send(d)
        }
    }
}

// handleRecevier -
func (c *Connection) handleRecevier() error {
    defer c.bzn.HandleConnClose(c)
    // defer c.Close()
    // buffer := make([]byte, c.maxBufferLen)
    // currOffset := 0
    // readOffset := 0
    // remainLen := 0
    // nextRead := true
    count := 0
    for {
        // if nextRead {
        //     len, err := c.Read(buffer[readOffset:])
        //     if err != nil {
        //         logbz.Errorf("Connection handleRecevier - read error.", err.Error())
        //         return err
        //     }

        //     readOffset += len
        //     remainLen = readOffset - currOffset
        //     if readOffset > c.maxBufferLen - c.maxBufferLen / 10 {
        //         buffer = buffer[currOffset:readOffset]
        //         currOffset = 0; readOffset = remainLen
        //     }
        // }
        if _, err := c.BufferRead.Read(func(b []byte) (int, error) {
            return c.TCPConn.Read(b)
        }); err != nil {
            logbz.Errorf("Connection handleRecevier - read error.", err.Error())
            return err
        }

        // len, err := c.TCPConn.Read(buffer)
        // if err != nil {
        //     logbz.Errorf("Connection handleRecevier - read error.", err.Error())
        //     return err
        // }

        if c.BufferRead.Empty() {
            logbz.Debugln("Connection handleRecevier - wait next.")
            continue
        }

        out := &protocol.CommonPt{}
        // if err := protocol.Unmarshal(buffer[0:len], out); err != nil {
        if err := protocol.Unmarshal(c.BufferRead.Get(), out); err != nil {
                // if err := protocol.Unmarshal(buffer[currOffset:readOffset], out); err != nil {
            if err == protocol.ErrNoFixedMe {
                logbz.Errorln("Connection handleRecevier - Unmarshal error.", err.Error())
                return err
            }
            // logbz.Errorln("Connection handleRecevier - Unmarshal ------- error.", err.Error())
            c.BufferRead.Step()
            continue
        }
        c.BufferRead.Next(out.Len())

        // currOffset += out.Len()
        // remainLen = readOffset - currOffset

        // fmt.Printf(">>>>> Connection handleRecevier - recv buffer len %d, unmarshal %d, count %d, payload(%d).\n", c.BufferRead.Length(), out.Len(), count, out.Length)
        // fmt.Printf(">>>>> Connection handleRecevier - recv buffer len %d, unmarshal %d, count %d.\n", len, out.Len(), count)
        if out.Type == protocol.Method_CHANNEL_CREATE {
            c.durationMs = utils.NowDiff(int64(out.Timestamp)).Milliseconds()
        }
        ms := utils.Abs(utils.NowDiff(int64(out.Timestamp)).Milliseconds() - c.durationMs)
        if ms > 3000 {
            fmt.Printf("Connection.handleRecevier - out %v, dura: %d(%d) ms, ts: %v, %v\n", out, ms, c.durationMs,
                utils.MsFormat(int64(out.Timestamp)), utils.MsFormat(time.Now().UnixMilli()))
        }

        count++
        if err := c.bzn.HandlePt(c, out); err != nil {
            return err
        }

        if c.quit == true {
            break
        }
    }
    return nil
}


