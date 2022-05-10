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
    sendPeriod int // Microsecond
    sendBufferLen int // bytes.
    recvBufferLen int //


    // stats.
    sentCount int64 //
    sentBytes int64
    sendBeginTime time.Time
    sendEndTime time.Time
    recvCount int64
    recvBytes int64
    recvBeginTime time.Time
    recvEndTime time.Time

}

// NewAppsClient -
func NewAppsClient(tcpAddress string) *AppsClient {
    c := &AppsClient{
        done: make(chan bool),
        sendPeriod: sendPeroid,
        sendBufferLen: maxBufferLen,
        recvBufferLen: maxBufferLen,
        tcpAddress: tcpAddress,
    }
    return c
}

// show -
func (app *AppsClient) show() {
    fmt.Printf("send (time: %v => %v) - %v bytes(%v count)\n", app.sendBeginTime.Format("2006-01-02 15:04:05.999999999"), app.sendEndTime.Format("2006-01-02 15:04:05.999999999"), app.sentBytes, app.sentCount)
    fmt.Printf("recv (time: %v => %v) - %v bytes(%v count)\n", app.recvBeginTime.Format("2006-01-02 15:04:05.999999999"), app.recvEndTime.Format("2006-01-02 15:04:05.999999999"), app.recvBytes, app.recvCount)
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

// handleSender -
func (app *AppsClient) handleSender() {
    app.sendBeginTime = time.Now()
    // buffer := make([]byte, app.sendBufferLen)
    buffer := utils.RandomBytes(app.sendBufferLen, nil)
    sendDuration := time.Duration(app.sendPeriod) * time.Millisecond
    ticker := time.NewTicker(sendDuration)
    defer ticker.Stop()
    fmt.Printf("AppsClient.handleSender - send duration %d ms, buffer len %d, begin time %v.\n", app.sendPeriod, app.sendBufferLen, app.sendBeginTime.Format("2006-01-02 15:04:05"))
    bQuit := false
    for {
        select {
        case <- app.done:
            bQuit = true
        case <- ticker.C:
            dura := utils.NewDuration()
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
            fmt.Printf("send buffer No.%d, len %d, real %d. dura %d ms.\n", app.sentCount, app.sendBufferLen, n, dura.DuraMs())
            // fmt.Printf("send buffer No.%d, len %d, real %d. dura %d ms. =>%s.\n", app.sentCount, app.sendBufferLen, n, dura.DuraMs(), string(buffer[0:10]))
            app.sentCount += 1
            app.sentBytes += int64(n)
        }
        if bQuit == true {
            break
        }
    }
    app.sendEndTime = time.Now()
    fmt.Printf("AppsClient.handleSender - send duration %d ms, buffer len %d, end time %v ms(dura %v ms).\n", app.sendPeriod, app.sendBufferLen, app.sendEndTime.Format("2006-01-02 15:04:05"), app.sendEndTime.Sub(app.sendBeginTime))
}

// handleRecevicer -
func (app *AppsClient) handleRecevicer() {
    buffer := make([]byte, app.recvBufferLen)
    for {
        n, err := app.tcpClient.Read(buffer)
        if err != nil {
            fmt.Printf("AppsClient.handleRecevicer error.%v.\n", err.Error())
            break
        }
        if app.recvCount == 0 {
            app.recvBeginTime = time.Now()
            fmt.Printf("AppsClient.handleRecevicer - begin. recv begin %v.\n", app.recvBeginTime)
        }
        app.recvCount += 1
        app.recvBytes += int64(n)
        if n != app.recvBufferLen {
            fmt.Printf("AppsClient.handleRecevicer recv buffer len %d not equal send buffer, real %d.\n", app.recvBufferLen, n)
            break
        }
    }
    app.recvEndTime = time.Now()
    fmt.Printf("AppsClient.handleRecevicer - end... recv begin %v, end %v, dura %v.\n", app.recvBeginTime, app.recvEndTime, app.recvEndTime.Sub(app.recvBeginTime))
    app.done <- true
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
        } else {
            fmt.Printf("cmd => (%v) not impliment.\r\n", cmd)
        }
    }
    fmt.Printf("\nAppsClient.wait - end...\n")
    return nil
}


// Main -
func (app *AppsClient) Main() error {
    tcpClient := bzc.NewTcpClient(app.tcpAddress)
    if err := tcpClient.Dial(); err != nil {
        fmt.Println("AppsClient.Main error", err.Error())
        return err
    }
    app.tcpClient = tcpClient
    // go app.handleSender()
    // go app.handleRecevicer()
    return app.wait()
}
