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
}

var _ bz.BZNet = (*BytezeroNet)(nil)

// NewBytezeroNet -
func NewBytezeroNet(ctx context.Context, done chan bool) *BytezeroNet {
    bzn := &BytezeroNet{
        done: done,
    }
    return bzn
}

// Main -
func (bzn *BytezeroNet) Main() {
    logbz.Debugln("BytezeroNet Main...")
    go bzn.StartTcp()
}

// Quit -
func (bzn *BytezeroNet) Quit() bool {
    logbz.Debugln("BytezeroNet maybe quit...")
    return true
}

func (bzn *BytezeroNet) StartTcp() {
    tcpServer := server.NewTcpServer(":7788")
    err := tcpServer.Listen()
    if err != nil {
        logbz.Errorln("BytezeroNet.StartTcp.Listen error.%v.", err.Error())
        bzn.done <- true
    }
}

