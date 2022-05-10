package cores

import (
	"context"

	bz "github.com/alackfeng/bytezero/bytezero"
	"github.com/alackfeng/bytezero/cores/utils"
)

var logbz = utils.Logger(utils.Fields{"animal": "main"})

// BytezeroNet - BytezeroNet
type BytezeroNet struct {

}

var _ bz.BZNet = (*BytezeroNet)(nil)

// NewBytezeroNet -
func NewBytezeroNet(ctx context.Context, done chan bool) *BytezeroNet {
    bzn := &BytezeroNet{}
    return bzn
}

// Main -
func (bzn *BytezeroNet) Main() {
    logbz.Debugln("BytezeroNet Main...")
}

// Quit -
func (bzn *BytezeroNet) Quit() bool {
    logbz.Debugln("BytezeroNet maybe quit...")
    return true
}

