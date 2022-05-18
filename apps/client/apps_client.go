package client

import (
	"bufio"
	"fmt"
	"os"
	"time"

	bzc "github.com/alackfeng/bytezero/cores/client"
	"github.com/alackfeng/bytezero/cores/utils"
)

const maxBufferLen int = 1024 * 1024 * 1
const sendPeroid int = 1000 // ms.


// AppsClient - 测试客户端.
type AppsClient struct {
    done chan bool

    tcpAddress string
    tcpClient* bzc.TcpClient

    udpAddress string
    udpClient* bzc.UdpClient

    sendPeriod int // Microsecond
    sendBufferLen int // bytes.
    recvBufferLen int //
    recvCheck bool // false, close connection.


    // stats.
    sendStat utils.StatBandwidth
    recvStat utils.StatBandwidth

}

// NewAppsClient -
func NewAppsClient() *AppsClient {
    c := &AppsClient{
        done: make(chan bool),
        sendPeriod: sendPeroid,
        sendBufferLen: maxBufferLen,
        recvBufferLen: maxBufferLen,
        recvCheck: false,
    }
    return c
}

// show -
func (app *AppsClient) show() {
    fmt.Printf("send %v\n", app.sendStat)
    fmt.Printf("recv %v\n", app.recvStat)
}

// SetSendPeroid -
func (app *AppsClient) SetSendPeroid(peroidMs int) *AppsClient {
    app.sendPeriod = peroidMs
    return app
}

// SetMaxBufferLen -
func (app *AppsClient) SetMaxBufferLen(len int) *AppsClient {
    app.sendBufferLen = len
    app.recvBufferLen = len
    return app
}

// SetRecvCheck -
func (app *AppsClient) SetRecvCheck(check bool) *AppsClient {
    app.recvCheck = check
    return app
}

// SetTcpAddress -
func (app *AppsClient) SetTcpAddress(address string) *AppsClient {
    app.tcpAddress = address
    return app
}

// SetUdpAddress -
func (app *AppsClient) SetUdpAddress(address string) *AppsClient {
    app.udpAddress = address
    return app
}

// handleSender -
func (app *AppsClient) handleSender() {
    if app.tcpAddress == "" {
        return
    }
    app.sendStat.Begin()
    // buffer := make([]byte, app.sendBufferLen)
    buffer := utils.RandomBytes(app.sendBufferLen, nil)
    // sendDuration := time.Duration(app.sendPeriod) * time.Millisecond
    // ticker := time.NewTicker(sendDuration)
    // defer ticker.Stop()
    fmt.Printf("AppsClient.handleSender - send duration %d ms, buffer len %d, begin time %v.\n", app.sendPeriod, app.sendBufferLen, app.sendStat.InfoAll())
    bQuit := false
    for {
        select {
        case <- app.done:
            bQuit = true
        // case <- ticker.C:
        default:
            // dura := utils.NewDuration()
            n, err := app.tcpClient.Write(buffer)
            if err != nil {
                fmt.Printf("AppsClient.handleSender - send error.%v.\n", err.Error())
                break
            }
            if n != app.sendBufferLen {
                fmt.Printf("AppsClient.handleSender - send buffer len %d not equal real sent %d.\n", app.sendBufferLen, n)
                break
            }
            // fmt.Printf("send buffer No.%d, len %d, real %d. =>%v.\n", app.sentCount, app.sendBufferLen, n, buffer[0:10])
            // fmt.Printf("send buffer No.%d, len %d, real %d.\n", app.sentCount, app.sendBufferLen, n)
            // if app.sendStat.Count % 1000 == 0 {
            //     fmt.Printf("send buffer No.%d, len %d, real %d. dura %d ms.\n", app.sendStat.Count, app.sendBufferLen, n, dura.DuraMs())
            // }
            // fmt.Printf("send buffer No.%d, len %d, real %d. dura %d ms. =>%s.\n", app.sentCount, app.sendBufferLen, n, dura.DuraMs(), string(buffer[0:10]))
            app.sendStat.Inc(int64(n))
            if app.sendStat.Count % 100 == 0 {
                time.Sleep(time.Millisecond*10)
            }
        }
        if bQuit == true {
            break
        }
    }
    app.sendStat.End()
    fmt.Printf("AppsClient.handleSender - send duration %d ms, buffer len %d, %v.\n", app.sendPeriod, app.sendBufferLen, app.sendStat.InfoAll())
}

// handleRecevicer -
func (app *AppsClient) handleRecevicer() {
    if app.tcpAddress == "" {
        return
    }
    buffer := make([]byte, app.recvBufferLen)
    currTime := time.Now()
    for {
        n, err := app.tcpClient.Read(buffer)
        if err != nil {
            fmt.Printf("AppsClient.handleRecevicer error.%v.\n", err.Error())
            break
        }
        if app.recvStat.Bytes == 0 {
            app.recvStat.Begin()
            fmt.Printf("AppsClient.handleRecevicer - begin. recv begin %v.\n", app.recvStat.Info())
        }
        app.recvStat.Inc(int64(n))
        if n != app.recvBufferLen {
            // fmt.Printf("AppsClient.handleRecevicer recv buffer len %d not equal send buffer, real %d.\n", app.recvBufferLen, n)
            if app.recvCheck {
                break
            }
        }
        if time.Now().Sub(currTime).Milliseconds() > 1000 {
            currTime = time.Now()
            fmt.Printf("AppsClient.handleRecevicer recv - count %d, bps %s. send bps %s\n", app.recvStat.Count, utils.ByteSizeFormat(app.recvStat.Bps1s()), utils.ByteSizeFormat(app.sendStat.Bps1s()))
        }
    }
    app.recvStat.End()
    fmt.Printf("AppsClient.handleRecevicer - end... %v.\n", app.recvStat.InfoAll())
    app.done <- true
}

// handleUdpRecevicer -
func (app *AppsClient) handleUdpRecevicer() {
    if app.udpAddress == "" {
        return
    }
    buffer := make([]byte, app.recvBufferLen)
    currTime := time.Now()
    for {
        n, addr, err := app.udpClient.ReadFromUDP(buffer[:])
        if err != nil {
            fmt.Printf("AppsClient.handleUdpRecevicer error.%v.\n", err.Error())
            break
        }
        if app.recvStat.Bytes == 0 {
            app.recvStat.Begin()
            fmt.Printf("AppsClient.handleUdpRecevicer - begin. recv begin %v - remote addr %v.\n", app.recvStat.Info(), addr)
        }
        app.recvStat.Inc(int64(n))
        if n != app.recvBufferLen {
            // fmt.Printf("AppsClient.handleUdpRecevicer recv buffer len %d not equal send buffer, real %d.\n", app.recvBufferLen, n)
            if app.recvCheck {
                break
            }
        }
        if time.Now().Sub(currTime).Milliseconds() > 1000 {
            currTime = time.Now()
            fmt.Printf("AppsClient.handleUdpRecevicer recv - count %d, bps %s. send bps %s.\n", app.recvStat.Count, utils.ByteSizeFormat(app.recvStat.Bps1s()), utils.ByteSizeFormat(app.sendStat.Bps1s()))
        }
    }
    app.recvStat.End()
    fmt.Printf("AppsClient.handleUdpRecevicer - end... %v.\n", app.recvStat.InfoAll())
    app.done <- true
}

// handleUdpSender -
func (app *AppsClient) handleUdpSender() {
    if app.udpAddress == "" {
        return
    }
    app.sendStat.Begin()
    buffer := utils.RandomBytes(app.sendBufferLen, nil)
    sendDuration := time.Duration(app.sendPeriod) * time.Millisecond
    ticker := time.NewTicker(sendDuration)
    defer ticker.Stop()
    fmt.Printf("AppsClient.handleUdpSender - send duration %d ms, buffer len %d, begin time %v.\n", app.sendPeriod, app.sendBufferLen, app.sendStat.InfoAll())
    bQuit := false
    for {
        select {
        case <- app.done:
            bQuit = true
        case <- ticker.C:
        // default:
            dura := utils.NewDuration()
            n, err := app.udpClient.Write(buffer)
            if err != nil {
                fmt.Printf("AppsClient.handleUdpSender - send error.%v.\n", err.Error())
                break
            }
            if n != app.sendBufferLen {
                fmt.Printf("AppsClient.handleUdpSender - send buffer len %d not equal real sent %d.\n", app.sendBufferLen, n)
                break
            }
            // fmt.Printf("send buffer No.%d, len %d, real %d. =>%v.\n", app.sentCount, app.sendBufferLen, n, buffer[0:10])
            // fmt.Printf("send buffer No.%d, len %d, real %d.\n", app.sentCount, app.sendBufferLen, n)
            if app.sendStat.Count % 1000 == 0 {
                fmt.Printf("send buffer No.%d, len %d, real %d. dura %d ms.\n", app.sendStat.Count, app.sendBufferLen, n, dura.DuraMs())
            }
            // fmt.Printf("send buffer No.%d, len %d, real %d. dura %d ms. =>%s.\n", app.sentCount, app.sendBufferLen, n, dura.DuraMs(), string(buffer[0:10]))
            app.sendStat.Inc(int64(n))
        }
        if bQuit == true {
            break
        }
    }
    app.sendStat.End()
    fmt.Printf("AppsClient.handleUdpSender - send duration %d ms, buffer len %d, %v.\n", app.sendPeriod, app.sendBufferLen, app.sendStat.InfoAll())
}

// wait -
func (app *AppsClient) wait() error {
    fmt.Printf("AppsClient.wait - begin.\n")
    input := bufio.NewScanner(os.Stdin)
    for {
        fmt.Printf("AppsClient - cmd: ")
        if input.Scan() == false {
            break
        }
        cmd := input.Text()
        if cmd == "" {
        } else if cmd == "quit" || cmd == "q" || cmd == "Q" {
            break
        } else if cmd == "info" || cmd == "i" {
            app.show()
        } else if cmd == "send" {
            go app.handleRecevicer()
            go app.handleSender()
        } else if cmd == "udp" {
            go app.handleUdpRecevicer()
            go app.handleUdpSender()
        } else {
            fmt.Printf("cmd => (%v) not impliment.\r\n", cmd)
        }
    }
    fmt.Printf("\nAppsClient.wait - end...\n")
    return nil
}


// Main -
func (app *AppsClient) Main() error {
    if app.tcpAddress != "" {
        tcpClient := bzc.NewTcpClient(app.tcpAddress)
        if err := tcpClient.Dial(); err != nil {
            fmt.Println("AppsClient.Main tcp error", err.Error())
            return err
        }
        app.tcpClient = tcpClient
    }
    if app.udpAddress != "" {
        udpClient := bzc.NewUdpClient(app.udpAddress)
        if err := udpClient.Dial(); err != nil {
            fmt.Println("AppsClient.Main udp error", err.Error())
            return err
        }
        app.udpClient = udpClient
    }
    // go app.handleSender()
    // go app.handleRecevicer()
    return app.wait()
}
