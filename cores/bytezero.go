package cores

import (
	"context"
	"fmt"
	"net"

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
    go bzn.StartUdp()
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

// HandleConn -
func (bzn *BytezeroNet) HandleConn(tcpConn *net.TCPConn) error {
    conn := NewConnection(bzn, tcpConn).Main()
    remoteId := conn.RemoteAddr().String()
    bzn.connections[remoteId] = conn
    logbz.Infof("BytezeroNet HandleConn - create connection client remote<%s>, local<%s>.", remoteId, conn.LocalAddr())
    return nil
}


// HandlePt -
func (bzn *BytezeroNet) HandlePt(conn bz.BZNetReceiver, commonPt *protocol.CommonPt) error {
    switch commonPt.Type {
    case protocol.Method_CHANNEL_CREATE:
        channelCreatePb := protocol.NewChannelCreatePb()
        if err := channelCreatePb.Unmarshal(commonPt.Payload); err != nil {
            return fmt.Errorf("ChannelCreatePb Unmarshal error.%v", err.Error())
        }
        if c, ok := conn.(*Connection); ok {
            // Update DevcieId etc.
            c.Set(channelCreatePb)

            sessionId := string(channelCreatePb.SessionId)
            if channel, ok := bzn.channels[sessionId]; ok {
                channel.Join(c).Ack()
            } else {
                bzn.channels[sessionId] = NewChannel().Create(c)
            }
        }
    case protocol.Method_STREAM_CREATE:
        fallthrough
    case protocol.Method_STREAM_ACK:
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




