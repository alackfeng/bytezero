package cores

import (
	"context"
	"fmt"
	"net"
	"regexp"
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
    tl* server.TlsServer

    // web api.
    gw *web.GinWeb

    l sync.Mutex
    connections map[string]*Connection
    channels map[string]*Channel

    accessIpsAllow utils.AccessIpsAllow
    accessIpsDeny utils.AccessIpsDeny
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

// CredentialUrls -
func (bzn *BytezeroNet) CredentialUrls() []string {
    return ConfigGlobal().App.Credential.Urls
}



// MargicV - MARGIC_SHIFT for transport secret.
func (bzn *BytezeroNet) MargicV() (byte, bool) {
    config := ConfigGlobal()
    return byte(bz.MargicValue), config.App.Server.Margic
}

// Main -
func (bzn *BytezeroNet) Main() {
    logbz.Debugln("BytezeroNet Main...")
    go bzn.StartTcp()
    go bzn.StartWeb()
    go bzn.StartTls()
}

// Quit -
func (bzn *BytezeroNet) Quit() bool {
    logbz.Debugln("BytezeroNet maybe quit...")
    return true
}

// StartWeb -
func (bzn *BytezeroNet) StartWeb() {
    config := ConfigGlobal()
    if !config.App.Web.Http.UP && !config.App.Web.Https.UP {
        return
    }
    if bzn.gw == nil {
        bzn.gw = web.NewGinWeb(config.App.Web.Http.Host, config.App.Web.Http.Heart, bzn)
    }
    if config.App.Web.Https.UP {
        bzn.gw.SetSecretTransport(config.App.Web.Https.Host, config.App.Web.Https.CaCert, config.App.Web.Https.CaKey)
    }
    bzn.gw.SetStaticInfo(config.App.Web.Static.Memory, config.App.Web.Static.LogPath, config.App.Web.Static.UploadPath)
    bzn.gw.Start()
}
// StartTcp -
func (bzn *BytezeroNet) StartTcp() {
    config := ConfigGlobal()
    if !config.App.Server.UP {
        return
    }
    tcpServer := server.NewTcpServer(bzn, config.App.Server.Address(), config.App.MaxBufferLen, config.App.RWBufferLen)
    err := tcpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartTcp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.ts = tcpServer
}

// StartTls -
func (bzn *BytezeroNet) StartTls() {
    config := ConfigGlobal()
    if !config.App.Tls.UP {
        return
    }
    tlsServer := server.NewTlsServer(bzn, config.App.Tls.Address(), config.App.Tls.CaCert, config.App.Tls.CaKey)
    err := tlsServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartTls.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.tl = tlsServer
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

// AccessIpsAllow -
func (bzn *BytezeroNet) AccessIpsAllow(ip string) error {
    config := ConfigGlobal()
    if config.App.AccessIpsAllow == "" {
        return nil
    }
    if len(bzn.accessIpsAllow.Allow) == 0 {
        if err := bzn.accessIpsAllow.Load(config.App.AccessIpsAllow); err != nil {
            return err
        }
        fmt.Println("-----------------AccessIpsAllow accessIps: ", bzn.accessIpsAllow.Allow)
        if len(bzn.accessIpsAllow.Allow) == 0 {
            return nil // allow all.
        }
    }

    if _, ok := bzn.accessIpsAllow.Allow[ip]; ok {
        return nil // allow.
    }
    match := false
    for ipAllow, drop := range bzn.accessIpsAllow.Allow {
        if ipAllow == "" || !drop {
            continue
        }
        if ipAllow == "*" {
            match = true
            break
        }
        if ipAllow == ip {
            match = true
            break
        }
        ipAllow = "^" + ipAllow + "$"
        if ok, _ := regexp.MatchString(ipAllow, ip); ok {
			match = true
			break
		}
    }
    if match {
        return nil // allow.
    }
    return fmt.Errorf("ip<%s> deny", ip)
}

// AccessIpsDeny -
func (bzn *BytezeroNet) AccessIpsDeny(ip string) error {
    config := ConfigGlobal()
    if config.App.AccessIpsDeny == "" {
        return nil
    }
    if len(bzn.accessIpsDeny.Deny) == 0 {
        if err := bzn.accessIpsDeny.Load(config.App.AccessIpsDeny); err != nil {
            return err
        }
        fmt.Println("-----------------AccessIpsDeny accessIps: ", bzn.accessIpsDeny.Deny)
        if len(bzn.accessIpsDeny.Deny) == 0 {
            return nil // allow all.
        }
    }

    if deny, ok := bzn.accessIpsDeny.Deny[ip]; ok && deny {
        return fmt.Errorf("ip<%s> deny", ip)
    }
    return nil // allow.
}

// AccessIpsForbid -
func (bzn *BytezeroNet) AccessIpsForbid(ip string, deny bool) error {
    config := ConfigGlobal()
    return bzn.accessIpsDeny.Upload(config.App.AccessIpsDeny, ip, deny)
}

// AccessIpsReload -
func (bzn *BytezeroNet) AccessIpsReload(allow bool) error {
    config := ConfigGlobal()
    if allow {
        if err := bzn.accessIpsAllow.Load(config.App.AccessIpsAllow); err != nil {
            return err
        }
        fmt.Println("-----------------AccessIpsReload allow accessIps: ", bzn.accessIpsDeny.Deny)
    } else {
        if err := bzn.accessIpsDeny.Load(config.App.AccessIpsDeny); err != nil {
            return err
        }
        fmt.Println("-----------------AccessIpsReload deny accessIps: ", bzn.accessIpsDeny.Deny)
    }
    return nil
}

// HandleConn -
func (bzn *BytezeroNet) HandleConn(conn net.Conn) error {
    // access it?
    if host, _, err := net.SplitHostPort(conn.RemoteAddr().String()); err == nil {
        if err := bzn.AccessIpsDeny(host); err != nil {
            conn.Close()
            return fmt.Errorf("ip<%s> deny", host)
        }
    }

    //
    bzn.l.Lock()
    c := NewConnection(bzn, conn).Main()
    bzn.connections[c.Id()] = c
    bzn.l.Unlock()
    logbz.Infof("BytezeroNet HandleConn - create connection id:<%s>.", c.Id())
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
                bzn.l.Unlock()
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
            bzn.l.Lock()
            channId := c.ChannId()
            if channel, ok := bzn.channels[channId]; ok {
		/* if commonPt.Type == protocol.Method_STREAM_DATA {
                    streamDataPb := &protocol.StreamDataPt{}
                    if err := commonPt.UnmarshalP(streamDataPb); err != nil {
                        fmt.Println("------- streamDataPb ", err.Error())
                    }
                    l := len(streamDataPb.Data) - 32
                    fmd5 := utils.GenMD5(string(streamDataPb.Data[0:l]))
                    fmd52 := string(streamDataPb.Data[l:len(streamDataPb.Data)])
                    if fmd5 != fmd52 {
                        fmt.Println("-------------error md5 ", streamDataPb.Total, streamDataPb.Offset, streamDataPb.Timestamp, fmd5, fmd52)
                    } else {
                        // fmt.Println("-------------md5 ", streamDataPb.Total, streamDataPb.Offset, streamDataPb.Timestamp, fmd5, fmd52)
                    }
                } */
                buf, err := commonPt.Marshal()
                if err != nil {
                    bzn.l.Unlock()
                    break
                }
                channel.Transit(func(c1, c2 *Connection) error {
                    if c == c1 {
                        return c2.Transit(buf)
                    }
                    return c1.Transit(buf)
                })
            }
            bzn.l.Unlock()
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




