package client

import (
	"fmt"
	"sync"
)

// AppsChannels -
type AppsChannels struct {
    app Client
    l sync.Mutex
    m map[string]*AppsChannel
}
var _ BaseHandle = (*AppsChannels)(nil)

// NewAppsChannels -
func NewAppsChannels(app Client) *AppsChannels {
    return &AppsChannels{
        app: app,
        m: make(map[string]*AppsChannel),
    }
}

// ConnectionChannel - BaseHandle interface.
func (a *AppsChannels) ConnectionChannel(sessionId string, observer ChannelObserver) (ChannelHandle, error) {
    a.l.Lock()
    defer a.l.Unlock()

    if channel, ok := a.m[sessionId]; ok {
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
    channel := NewAppsChannel(a.app)
    channel.RegisterObserver(observer)
    if err := channel.Start(a.app.TargetAddress(), sessionId); err != nil {
        return nil, fmt.Errorf("Connection Channel Start error.%v", err.Error())
    }
    a.m[sessionId] = channel
    return channel, nil
}

// ChannelClose - BaseHandle interface.
func (a *AppsChannels) ChannelClose(sessionId string) (err error) {
    a.l.Lock()
    if _, ok := a.m[sessionId]; ok {
        // err = channel.Stop()
        delete(a.m, sessionId)
    }
    a.l.Unlock()
    return nil
}

// ChannelCloseByHandle - BaseHandle interface.
func (a *AppsChannels) ChannelCloseByHandle(handle ChannelHandle) (err error) {
    a.l.Lock()
    if h, ok := handle.(*AppsChannel); ok {
        if channel, ok := a.m[h.sessionId]; ok {
            err = channel.Stop()
            delete(a.m, h.sessionId)
        }
    }
    a.l.Unlock()
    return nil
}
