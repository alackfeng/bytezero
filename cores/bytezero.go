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
}

var _ bz.BZNet = (*BytezeroNet)(nil)

// NewBytezeroNet -
func NewBytezeroNet(ctx context.Context, done chan bool) *BytezeroNet {
    bzn := &BytezeroNet{
        done: done,
        tsAddr: ":7788",
        usAddr: ":7789",
    }
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
    tcpServer := server.NewTcpServer(bzn.tsAddr)
    err := tcpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartTcp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.ts = tcpServer
}

func (bzn *BytezeroNet) StartUdp() {
    udpServer := server.NewUdpServer(bzn.usAddr)
    err := udpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartUdp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
    bzn.us = udpServer
}



