package cores

import (
	"context"
	"fmt"
	"net"
	"sync"

	bz "github.com/alackfeng/bytezero/bytezero"
	"github.com/alackfeng/bytezero/bytezero/protocol"
	"github.com/alackfeng/bytezero/cores/server"
	"github.com/alackfeng/bytezero/cores/utils"
)

var logbz = utils.Logger(utils.Fields{"animal": "main"})

// BytezeroNet - BytezeroNet
type BytezeroNet struct {
    done chan bool
    ts* server.TcpServer
    tsAddr string
    us* server.UdpServer
    usAddr string
    maxBufferLen int
    rwBufferLen int

    l sync.Mutex
    connections map[string]*Connection
    channels map[string]*Channel
}

var _ bz.BZNet = (*BytezeroNet)(nil)

// NewBytezeroNet -
func NewBytezeroNet(ctx context.Context, done chan bool) *BytezeroNet {
    bzn := &BytezeroNet{
        done: done,
        tsAddr: ":7788",
        usAddr: ":7789",
        maxBufferLen: 1024*1024*10,
        rwBufferLen: 1024,
        connections: make(map[string]*Connection),
        channels: make(map[string]*Channel),
    }
    return bzn
}

// SetMaxBufferLen -
func (bzn *BytezeroNet) SetPort(port int) *BytezeroNet {
    bzn.tsAddr = ":" + utils.IntToString(port)
    bzn.usAddr = ":" + utils.IntToString(port+1)
    return bzn
}

// SetMaxBufferLen -
func (bzn *BytezeroNet) SetMaxBufferLen(n int) *BytezeroNet {
    bzn.maxBufferLen = n
    return bzn
}

// SetRWBufferLen -
func (bzn *BytezeroNet) SetRWBufferLen(n int) *BytezeroNet {
    bzn.rwBufferLen = n
    return bzn
}

// Main -
func (bzn *BytezeroNet) Main() {
    logbz.Debugln("BytezeroNet Main...")
    go bzn.StartTcp()
    // go bzn.StartUdp()
}

// Quit -
func (bzn *BytezeroNet) Quit() bool {
    logbz.Debugln("BytezeroNet maybe quit...")
    return true
}

// StartTcp -
func (bzn *BytezeroNet) StartTcp() {
    tcpServer := server.NewTcpServer(bzn, bzn.tsAddr, bzn.maxBufferLen, bzn.rwBufferLen)
    err := tcpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartTcp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.ts = tcpServer
}

// HandleConnClose -
func (bzn *BytezeroNet) HandleConnClose(connection interface{}) {
    bzn.l.Lock()
    if c, ok := connection.(*Connection); ok {
        if channel, ok := bzn.channels[c.ChannId()]; ok {
            channel.LeaveAll()
            delete(bzn.channels, c.ChannId())
        }
        delete(bzn.connections, c.Id())
    }
    bzn.l.Unlock()
}

// HandleConn -
func (bzn *BytezeroNet) HandleConn(tcpConn *net.TCPConn) error {
    bzn.l.Lock()
    conn := NewConnection(bzn, tcpConn).Main()
    bzn.connections[conn.Id()] = conn
    bzn.l.Unlock()
    logbz.Infof("BytezeroNet HandleConn - create connection id:<%s>.", conn.Id())
    return nil
}


// HandlePt -
func (bzn *BytezeroNet) HandlePt(conn bz.BZNetReceiver, commonPt *protocol.CommonPt) error {
    switch commonPt.Type {
    case protocol.Method_CHANNEL_CREATE:
        channelCreatePb := protocol.NewChannelCreatePb()
        if err := commonPt.UnmarshalP(channelCreatePb); err != nil {
            return fmt.Errorf("ChannelCreatePb Unmarshal error.%v", err.Error())
        }
        fmt.Println("BytezeroNet.HandlePt - ", channelCreatePb)
        if c, ok := conn.(*Connection); ok {
            bzn.l.Lock()
            // Update DevcieId etc.
            if err := c.Set(channelCreatePb).Check(); err != nil {
                logbz.Errorf("BytezeroNet.HandlePt - connection check error.", err.Error())
                break
            }
            if channel, ok := bzn.channels[c.ChannId()]; ok {
                channel.Join(c).Ack(protocol.ErrCode_success, "ok")
            } else {
                bzn.channels[c.ChannId()] = NewChannel().Create(c)
            }
            bzn.l.Unlock()
        }
    case protocol.Method_STREAM_CREATE:
        fallthrough
    case protocol.Method_STREAM_ACK:
        fallthrough
    case protocol.Method_STREAM_DATA:
        fallthrough
    case protocol.Method_STREAM_CLOSE:
        if c, ok := conn.(*Connection); ok {
            channId := c.ChannId()
            if channel, ok := bzn.channels[channId]; ok {
                buf, err := commonPt.Marshal()
                if err != nil {
                    break
                }
                channel.Transit(func(c1, c2 *Connection) error {
                    if c == c1 {
                        return c2.Transit(buf)
                    }
                    return c1.Transit(buf)
                })
            }
        }

    default:
        logbz.Errorln("BytezeroNet HandlePt - Method type %v no impliment", commonPt.Type)
        return protocol.ErrBZProtocol
    }
    return nil
}

// StartUdp -
func (bzn *BytezeroNet) StartUdp() {
    udpServer := server.NewUdpServer(bzn.usAddr, bzn.maxBufferLen, bzn.rwBufferLen)
    err := udpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartUdp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.us = udpServer
}




