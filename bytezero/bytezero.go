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
    HandleConnClose(connection interface{})
    HandlePt(BZNetReceiver, *protocol.CommonPt) error

    // credential.
    AppID() string
    AppKey() string
    CredentialExpiredMs() int64
    MargicV() byte // MARGIC_SHIFT for transport secret.
}



// BZNetServer -
type BZNetServer interface {
}

// BZNetClient -
type BZNetClient interface {

}
