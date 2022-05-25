package client

import (
	"fmt"

	"github.com/alackfeng/bytezero/bytezero/protocol"
)

// AppsUploadResource -
type AppsUploadResource struct {
    app Client
    sessionId string
    channelId protocol.ChannelId
    streamId protocol.StreamId
    channelHandle ChannelHandle

    filePath string
}
var _ StreamObserver = (*AppsUploadResource)(nil)

// NewAppsUploadResource -
func NewAppsUploadResource(app Client, sessionId, filePath string) *AppsUploadResource {
    return &AppsUploadResource{
        app: app,
        sessionId: sessionId,
        filePath: filePath,
    }
}

func (a *AppsUploadResource) uploadFile() error {
    return nil
}

// Start -
func (a *AppsUploadResource) Start() error {
    // 创建Channel通道.
    channelHandle, err := a.app.ConnectionChannel(a.sessionId)
    if err != nil {
        return err
    }
    a.channelHandle = channelHandle

    // 创建Stream通道.
    channelHandle.StreamCreate(a.streamId, a)

    // 开始执行任务.
    return a.uploadFile()
}

// Stop -
func (a *AppsUploadResource) Stop() error {
    return nil
}

// OnStreamSuccess -
func (a *AppsUploadResource) OnStreamSuccess(sid protocol.StreamId) {
    fmt.Println("AppsUploadResource.OnStreamSuccess - stream id ", sid)
}

// OnStreamError -
func (a *AppsUploadResource) OnStreamError(code int, message string) {
    fmt.Printf("AppsUploadResource.OnStreamError - code %d, %s.\n", code, message)
}
