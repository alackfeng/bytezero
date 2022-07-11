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
const transDataLenghtDefault = 10240

// Connection -
type Connection struct {
    utils.BufferRead
    net.Conn
    bzn bz.BZNet
    quit bool
    durationMs int64

    // Info.
    DeviceId string
    SessionId string

    dataSync bool
    data chan []byte
}

var _ bz.BZNetReceiver = (*Connection)(nil)

// NewConnection -
func NewConnection(bzn bz.BZNet, c net.Conn) *Connection {
    cc := &Connection {
        Conn: c,
        bzn: bzn,
        BufferRead: *utils.NewBufferRead(defaultMaxBufferLen*1000),
        dataSync: true,
    }
    if !cc.dataSync {
        cc.data = make(chan []byte, transDataLenghtDefault)
    }
    // MARGIC_SHIFT for transport secret.
    cc.BufferRead.Margic, cc.BufferRead.Secret = bzn.MargicV()
    return cc
}

// Main -
func (c *Connection) Main()*Connection {
    if !c.dataSync {
        go c.handleTrans()
    }
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
    go (func() { c.Close() })()
    c.quit = true
}

func (c Connection) String() string {
    return fmt.Sprintf("Connection#[SessionId<%s>, DeviceId<%s>]", c.SessionId, c.DeviceId)
}


// Transit - to connection.
func (c *Connection) Transit(buf []byte) error {
    if !c.dataSync {
        c.data <- buf
    } else {
        return c.Send(buf)
    }
    return nil
}

// Send -
func (c *Connection) Send(buf []byte) error {
    // MARGIC_SHIFT for transport secret.
    if c.BufferRead.Secret {
        for i:=0; i<len(buf); i++ {
            buf[i] ^= c.BufferRead.Margic
        }
    }

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
            l := len(c.data)
            for i:=0; i<l; i++ {
                d, ok := <- c.data
                if !ok {
                    return
                }
                c.Send(d)
            }
        }
    }
}

// handleRecevier -
func (c *Connection) handleRecevier() error {
    defer c.bzn.HandleConnClose(c)

    for {

        if c.quit == true {
            break
        }

        if _, err := c.BufferRead.Read(func(b []byte) (int, error) {
            return c.Conn.Read(b)
        }); err != nil {
            logbz.Errorf("Connection handleRecevier - read error.", err.Error())
            return err
        }

        if c.BufferRead.Empty() {
            logbz.Debugln("Connection handleRecevier - wait next.")
            continue
        }

        out := &protocol.CommonPt{}
        if err := protocol.Unmarshal(c.BufferRead.Get(), out); err != nil {
            if err == protocol.ErrNoFixedMe || err == protocol.ErrNoMethodType {
                logbz.Errorln("Connection handleRecevier - Unmarshal error.", err.Error())
                if host, _, err := net.SplitHostPort(c.RemoteAddr().String()); err == nil {
                    c.bzn.AccessIpsForbid(host, false)
                }
                return err
            }
            // logbz.Errorln("Connection handleRecevier - Unmarshal ------- error.", err.Error())
            c.BufferRead.Step()
            continue
        }
        if out.Ver < protocol.CurrentVersion { // no support Version.
            logbz.Errorln("Connection handleRecevier - no support version.", out.Ver)
            return protocol.ErrNoSupportVersion
        }
        c.BufferRead.Next(out.Len())

        // fmt.Printf(">>>>> Connection handleRecevier - recv buffer len %d, unmarshal %d, count %d, payload(%d).\n", c.BufferRead.Length(), out.Len(), count, out.Length)
        // fmt.Printf(">>>>> Connection handleRecevier - recv buffer len %d, unmarshal %d, count %d.\n", len, out.Len(), count)
        if out.Type == protocol.Method_CHANNEL_CREATE {
            c.durationMs = utils.NowDiff(int64(out.Timestamp)).Milliseconds()
        }
        ms := utils.Abs(utils.NowDiff(int64(out.Timestamp)).Milliseconds() - c.durationMs)
        if ms > 3000 {
            fmt.Printf("Connection.handleRecevier - out %v, dura: %d(%d) ms, ts: %v, %v\n", out.Type, ms, c.durationMs,
                utils.MsFormat(int64(out.Timestamp)), utils.MsFormat(time.Now().UnixMilli()))
        }

        if err := c.bzn.HandlePt(c, out); err != nil {
            return err
        }

    }
    return nil
}


