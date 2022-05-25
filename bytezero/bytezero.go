package bytezero

import (
	"net"

	"github.com/alackfeng/bytezero/bytezero/protocol"
)

// BZNetReceiver -
type BZNetReceiver interface {
    ChannId() string
}

// BZNet -
type BZNet interface {
    HandleConn(*net.TCPConn) error
    HandlePt(BZNetReceiver, *protocol.CommonPt) error
}



// BZNetServer -
type BZNetServer interface {
}

// BZNetClient -
type BZNetClient interface {

}
