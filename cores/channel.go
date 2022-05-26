package cores

import (
	"fmt"

	"github.com/alackfeng/bytezero/bytezero/protocol"
)

// Channel -
type Channel struct {
    lc *Connection
    rt *Connection
    id protocol.ChannelId
}

// NewChannel -
func NewChannel() *Channel {
    return &Channel{}
}

// Online -
func (c *Channel) Online() bool {
    return !(c.lc == nil || c.rt == nil)
}

// Create -
func (c *Channel) Create(lc *Connection) *Channel {
    c.lc = lc
    c.id = protocol.GenChannelId()
    fmt.Println("Channel.Create - ", c.id)
    return c
}

// Join -
func (c *Channel) Join(o *Connection) *Channel {
    if c.rt == nil {
        c.rt = o
    } else if c.rt.Equals(o) { // Update, close last connection.
        c.rt.Quit()
        c.rt = o
    } else if c.lc == nil {
        c.lc = o
    } else if c.lc.Equals(o) {
        c.lc.Quit()
        c.lc = o
    }
    fmt.Println("Channel.Join - ", c.id, c.lc, c.rt)
    return c
}

// Ack -
func (c *Channel) Ack(code protocol.ErrCode, message string) error {
    if !c.Online() {
        return protocol.ErrNotBothOnline
    }
    channelAckPt := &protocol.ChannelAckPt{
        Id: c.id,
        Code: code,
        Message: []byte(message),
    }
    mByte, err := protocol.Marshal(channelAckPt)
    if err != nil {
        logbz.Errorf("Channel.Ack - error.%v.", err.Error())
        return err
    }
    fmt.Printf("Channel.Ack - send ack %v, to %v, %v.\n", channelAckPt, c.lc.Id(), c.lc.Id())
    if err := c.lc.Send(mByte); err != nil {
        logbz.Errorf("Channel.Ack - Send To %v, error.%v.", c.lc.Id(), err.Error())
        return err
    }
    if err := c.rt.Send(mByte); err != nil {
        logbz.Errorf("Channel.Ack - Send To %v, error.%v.", c.rt.Id(), err.Error())
        return err
    }
    return nil
}

// Transit -
func (c *Channel) Transit(send func(*Connection, *Connection) error)  {
    if send != nil {
        if err := send(c.lc, c.rt); err != nil {
            logbz.Errorln("Channel Transit send error", err.Error())
        }
    }
}


