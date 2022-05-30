package cores

import (
	"fmt"
	"net"

	bz "github.com/alackfeng/bytezero/bytezero"
	"github.com/alackfeng/bytezero/bytezero/protocol"
)

const defaultMaxBufferLen = 1538

// Connection -
type Connection struct {
    *net.TCPConn
    bzn bz.BZNet
    maxBufferLen int
    quit bool

    // Info.
    DeviceId string
    SessionId string
}

var _ bz.BZNetReceiver = (*Connection)(nil)

// NewConnection -
func NewConnection(bzn bz.BZNet, c *net.TCPConn) *Connection {
    return &Connection{
        TCPConn: c,
        bzn: bzn,
        maxBufferLen: defaultMaxBufferLen,
    }
}

// Main -
func (c *Connection) Main()*Connection {
    go c.handleSender()
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


// Transit - to connection.
func (c *Connection) Transit(buf []byte) error {
    return c.Send(buf)
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

// handleSender -
func (c *Connection) handleSender() {

}

// handleRecevier -
func (c *Connection) handleRecevier() error {
    defer c.Close()
    buffer := make([]byte, c.maxBufferLen)
    count := 0
    for {
        len, err := c.Read(buffer)
        if err != nil {
            logbz.Errorf("Connection handleRecevier - read error.", err.Error())
            return err
        }
        if len == 0 {
            logbz.Debugln("Connection handleRecevier - wait next.")
            continue
        }

        out := &protocol.CommonPt{}
        if err := protocol.Unmarshal(buffer[0:len], out); err != nil {
            logbz.Errorln("Connection handleRecevier - Unmarshal error.", err.Error())
            continue
        }

        fmt.Printf(">>>>> Connection handleRecevier - recv buffer len %d, unmarshal %d, count %d.\n", len, out.Len(), count)
        count++
        c.bzn.HandlePt(c, out)

        // wlen, err := c.Write(buffer[0:len])
        // if err != nil {
        //     return err
        // }
        // fmt.Printf("Connection handleRecevier - read %d, write %d.\n", len, wlen)

        if c.quit == true {
            break
        }
    }
    return nil
}


