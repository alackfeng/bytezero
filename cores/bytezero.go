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
	"github.com/alackfeng/bytezero/cores/web"
)

var logbz = utils.Logger(utils.Fields{"animal": "main"})

// BytezeroNet - BytezeroNet
type BytezeroNet struct {
    done chan bool
    ctx    context.Context
    appIds map[utils.OSType]string

    // server.
    ts* server.TcpServer
    us* server.UdpServer

    // web api.
    gw *web.GinWeb

    l sync.Mutex
    connections map[string]*Connection
    channels map[string]*Channel
}

var _ bz.BZNet = (*BytezeroNet)(nil)

// NewBytezeroNet -
func NewBytezeroNet(ctx context.Context, done chan bool) *BytezeroNet {
    bzn := &BytezeroNet{
        done: done,
        ctx: ctx,
        connections: make(map[string]*Connection),
        channels: make(map[string]*Channel),
    }
    bzn.appIds = map[utils.OSType]string{
        utils.OSTypeWindows: "xxx",
    }
    return bzn
}

// AppID -
func (bzn *BytezeroNet) AppID() string {
    return ConfigGlobal().App.Appid
}

// AppKey -
func (bzn *BytezeroNet) AppKey() string {
    return ConfigGlobal().App.Appkey
}

// CredentialExpiredMs -
func (bzn *BytezeroNet) CredentialExpiredMs() int64 {
    return ConfigGlobal().App.Credential.ExpiredMs
}
// Main -
func (bzn *BytezeroNet) Main() {
    logbz.Debugln("BytezeroNet Main...")
    go bzn.StartTcp()
    go bzn.StartWeb()
}

// Quit -
func (bzn *BytezeroNet) Quit() bool {
    logbz.Debugln("BytezeroNet maybe quit...")
    return true
}

// StartWeb -
func (bzn *BytezeroNet) StartWeb() {
    config := ConfigGlobal()
    if config.App.Web.Host == "" {
        return
    }
    if bzn.gw == nil {
        bzn.gw = web.NewGinWeb(config.App.Web.Host, config.App.Web.Heart, bzn)
    }
    bzn.gw.Start()
}
// StartTcp -
func (bzn *BytezeroNet) StartTcp() {
    config := ConfigGlobal()
    tcpServer := server.NewTcpServer(bzn, config.App.Server.Address(), config.App.MaxBufferLen, config.App.RWBufferLen)
    err := tcpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartTcp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.ts = tcpServer
}

// HandleConnClose -
func (bzn *BytezeroNet) HandleConnClose(connection interface{}) {
    fmt.Println("BytezeroNet::HandleConnClose - ")
    bzn.l.Lock()
    if c, ok := connection.(*Connection); ok {
        if channel, ok := bzn.channels[c.ChannId()]; ok {
            channel.LeaveAll()
            delete(bzn.channels, c.ChannId())
        }
        delete(bzn.connections, c.Id())
    }
    bzn.l.Unlock()
    fmt.Println("BytezeroNet::HandleConnClose - over.")
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
        channelCreatePb := &protocol.ChannelCreatePt{}
        if err := commonPt.UnmarshalP(channelCreatePb); err != nil {
            return fmt.Errorf("ChannelCreatePb Unmarshal error.%v", err.Error())
        }
        fmt.Println("BytezeroNet.HandlePt - ", channelCreatePb)
        if !utils.CheckAppID(utils.OSType(channelCreatePb.OS), string(channelCreatePb.AppId)) {
            logbz.Errorln("BytezeroNet.HandlePt - CheckAppID error.", channelCreatePb)
            return fmt.Errorf("OS not match AppID")
        }
        if err := utils.CredentialVerify(string(channelCreatePb.User), string(channelCreatePb.Sign), bzn.AppKey(), func(s string) ([]byte) {
            return channelCreatePb.FieldsSign()
        }); err != nil {
            logbz.Errorf("BytezeroNet.HandlePt - connection sign error.%s", err.Error())
            return err
        }

        if c, ok := conn.(*Connection); ok {
            bzn.l.Lock()
            // Update DevcieId etc.
            if err := c.Set(channelCreatePb).Check(); err != nil {
                logbz.Errorf("BytezeroNet.HandlePt - connection check error.%s", err.Error())
                break
            }
            if channel, ok := bzn.channels[c.ChannId()]; ok {
                channel.Join(c).Ack(protocol.ErrCode_success, "ok")
            } else {
                bzn.channels[c.ChannId()] = NewChannel().Create(c)
            }
            bzn.l.Unlock()
        }
    case protocol.Method_STREAM_DATA:
   /*     if c, ok := conn.(*Connection); ok {
            channId := c.ChannId()
            if _, ok := bzn.channels[channId]; ok {
                _, err := commonPt.Marshal()
                if err != nil {
                    break
                }
	    }
	}
        return nil
*/
	fallthrough
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
    config := ConfigGlobal()
    udpServer := server.NewUdpServer(config.App.Server.Address(), config.App.MaxBufferLen, config.App.RWBufferLen)
    err := udpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartUdp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.us = udpServer
}




