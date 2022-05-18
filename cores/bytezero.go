package cores

import (
	"context"

	bz "github.com/alackfeng/bytezero/bytezero"
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
    tcpServer := server.NewTcpServer(bzn.tsAddr, bzn.maxBufferLen, bzn.rwBufferLen)
    err := tcpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartTcp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.ts = tcpServer
}

func (bzn *BytezeroNet) StartUdp() {
    udpServer := server.NewUdpServer(bzn.usAddr, bzn.maxBufferLen, bzn.rwBufferLen)
    err := udpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartUdp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.us = udpServer
}



