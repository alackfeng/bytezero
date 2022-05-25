package client

import (
	"fmt"
	"sync"

	"github.com/alackfeng/bytezero/bytezero/protocol"
)

// AppsChannels -
type AppsChannels struct {
    app Client
    l sync.Mutex
    m map[string]*AppsChannel
}
var _ BaseHandle = (*AppsChannels)(nil)
var _ ChannelObserver = (*AppsChannels)(nil)

// NewAppsChannels -
func NewAppsChannels(app Client) *AppsChannels {
    return &AppsChannels{
        app: app,
        m: make(map[string]*AppsChannel),
    }
}

// ConnectionChannel - BaseHandle interface.
func (a *AppsChannels) ConnectionChannel(sessionId string) (ChannelHandle, error) {
    a.l.Lock()
    defer a.l.Unlock()
    channel, ok := a.m[sessionId]
    if ok {
        if channel.Online() {
            return channel, nil
        }
        channel.Stop()
        if err := channel.Start(channel.Address(), sessionId); err == nil {
            return channel, nil
        }
        // failed, remove it.
        delete(a.m, sessionId)
    }
    channel = NewAppsChannel(a.app)
    if err := channel.Start(a.app.TargetAddress(), sessionId); err != nil {
        return nil, fmt.Errorf("Connection Channel Start error.%v", err.Error())
    }
    a.m[sessionId] = channel
    return channel, nil
}

// ChannelClose - BaseHandle interface.
func (a *AppsChannels) ChannelClose(sessionId string) (err error) {
    a.l.Lock()
    if channel, ok := a.m[sessionId]; ok {
        err = channel.Stop()
        delete(a.m, sessionId)
    }
    a.l.Unlock()
    return nil
}

// OnSuccess - ChannelObserver interface.
func (a *AppsChannels) OnSuccess(cid protocol.ChannelId) {
    fmt.Println("AppsChannels.OnSuccess - channel id ", cid)
}

// OnError - ChannelObserver interface.
func (a *AppsChannels) OnError(code int, message string) {
    fmt.Printf("AppsChannels.OnError - code %d, %s.\n", code, message)
}

// onStream - ChannelObserver interface.
func (a *AppsChannels) onStream(s StreamHandle) (protocol.ErrCode, error) {
    fmt.Printf("AppsChannels.onStream - stream#%v.\n", s)
    return 0, nil
}
