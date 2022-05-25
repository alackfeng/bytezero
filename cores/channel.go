package cores

import "github.com/alackfeng/bytezero/bytezero/protocol"

// Channel -
type Channel struct {
    lc *Connection
    rt *Connection
    id string
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
    return c
}

// Join -
func (c *Channel) Join(o *Connection) *Channel {
    // if lc != nil
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
    return c
}

// Ack -
func (c *Channel) Ack() error {
    channelAckPt := protocol.NewChannelAckPt()
    mByte, err := protocol.Marshal(channelAckPt)
    if err != nil {
        return err
    }
    if !c.Online() {
        return protocol.ErrNotBothOnline
    }
    if n, err := c.lc.Write(mByte); err != nil {
        if n != len(mByte) {
            return protocol.ErrBufferNotAllWrite
        }
        return err
    }
    if n, err := c.rt.Write(mByte); err != nil {
        if n != len(mByte) {
            return protocol.ErrBufferNotAllWrite
        }
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


