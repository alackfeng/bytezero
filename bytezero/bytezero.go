package bytezero

import (
	"net"

	"github.com/alackfeng/bytezero/bytezero/protocol"
)

// MARGIC_SHIFT for transport secret.
const MargicValue = 0xA8

// BZNetReceiver -
type BZNetReceiver interface {
    ChannId() string
}

// BZNet -
type BZNet interface {
    HandleConn(net.Conn) error
    HandleConnClose(connection interface{})
    HandlePt(BZNetReceiver, *protocol.CommonPt) error

    // credential.
    AppID() string
    AppKey() string
    CredentialExpiredMs() int64
    CredentialUrls() []string
    MargicV() (byte, bool) // MARGIC_SHIFT for transport secret.
}



// BZNetServer -
type BZNetServer interface {
}

// BZNetClient -
type BZNetClient interface {

}
