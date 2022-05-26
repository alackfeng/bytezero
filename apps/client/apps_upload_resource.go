package client

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alackfeng/bytezero/bytezero/protocol"
)

type Mode uint8
const (
    ModeNone Mode = iota
    ModeUpload
    ModeAnswer
    ModeDownload
    ModeOffer
)

// AppsUploadResource -
type AppsUploadResource struct {
    app Client
    sessionId string
    channelId protocol.ChannelId
    streamId protocol.StreamId

    channelHandle ChannelHandle
    streamHandle StreamHandle

    filePath string
    bufferLen int
    mode Mode
}
var _ StreamObserver = (*AppsUploadResource)(nil)
var _ ChannelObserver = (*AppsUploadResource)(nil)

// NewAppsUploadResourceUpload -
func NewAppsUploadResourceUpload(app Client, sessionId, filePath string, bufferLen int) *AppsUploadResource {
    return &AppsUploadResource{
        app: app,
        sessionId: sessionId,
        filePath: filePath,
        bufferLen: bufferLen,
        mode: ModeUpload,
    }
}

// NewAppsUploadResourceAnswer -
func NewAppsUploadResourceAnswer(app Client, sessionId string) *AppsUploadResource {
    return &AppsUploadResource{
        app: app,
        sessionId: sessionId,
        mode: ModeAnswer,
    }
}

// uploadFile -
func (a *AppsUploadResource) uploadFile() (err error) {
    if a.mode != ModeUpload {
        fmt.Printf("AppsUploadResource.uploadFile - mode.%v not upload.\n", a.mode)
        return nil
    }
    if a.channelHandle == nil {
        fmt.Printf("AppsUploadResource.uploadFile - channelHandle is null.\n")
        return nil
    }

    // 创建Stream通道.
    a.streamHandle, err = a.channelHandle.StreamCreate(a.streamId, a)
    if err != nil {
        fmt.Printf("AppsUploadResource.uploadFile - streamHandle is null, error.%v.\n", err.Error())
        return err
    }

    fmt.Printf("AppsUploadResource.uploadFile - begin upload file<%s>, at Channel#%dStream#%d.\n", a.filePath, a.channelHandle.Id(), a.streamHandle.StreamId())
    f, err := os.Open(a.filePath)
    if err != nil {
        return err
    }
    defer f.Close()
    buf := make([]byte, a.bufferLen)
    for {
        n, err := f.Read(buf)
        if err != nil {
            if err == io.EOF {
                // return fmt.Errorf("File EOF")
                fmt.Printf("AppsUploadResource.uploadFile - upload file<%s> EOF.\n", a.filePath)
                break
            }
            return err
        }
        fmt.Printf("AppsUploadResource.uploadFile - send to Channel#%dStream#%d, buffer %d.\n", a.channelHandle.Id(), a.streamHandle.StreamId(), n)
        if err := a.streamHandle.SendData(buf[0:n]); err != nil {
            return err
        }
        time.Sleep(time.Millisecond * 1000)
    }
    fmt.Printf("AppsUploadResource.uploadFile - begin upload file<%s>, at Channel#%dStream#%d over.\n", a.filePath, a.channelHandle.Id(), a.streamHandle.StreamId())
    return nil
}

// Start -
func (a *AppsUploadResource) Start() (err error) {
    // 创建Channel通道.
    a.channelHandle, err = a.app.ConnectionChannel(a.sessionId, a)
    if err != nil {
        return err
    }

    // 开始执行任务.
    go a.uploadFile()
    return nil
}

// Stop -
func (a *AppsUploadResource) Stop() error {
    return nil
}


// OnSuccess - ChannelObserver interface.
func (a *AppsUploadResource) OnSuccess(cid protocol.ChannelId) {
    fmt.Printf("AppsUploadResource.OnSuccess - Channel#%d\n", cid)
}

// OnError - ChannelObserver interface.
func (a *AppsUploadResource) OnError(code int, message string) {
    fmt.Printf("AppsUploadResource.OnError - code %d, %s.\n", code, message)
}

// onStream - ChannelObserver interface.
func (a *AppsUploadResource) onStream(s StreamHandle) (protocol.ErrCode, error) {
    fmt.Printf("AppsUploadResource.onStream - Stream#%d\n", s.StreamId())
    a.streamHandle = s
    a.streamHandle.RegisterObserver(a)
    return protocol.ErrCode_success, fmt.Errorf("ok")
}


// OnStreamSuccess -
func (a *AppsUploadResource) OnStreamSuccess(sid protocol.StreamId) {
    fmt.Println("AppsUploadResource.OnStreamSuccess - stream id ", sid)
}

// OnStreamError -
func (a *AppsUploadResource) OnStreamError(code int, message string) {
    fmt.Printf("AppsUploadResource.OnStreamError - code %d, %s.\n", code, message)
}

// OnStreamData -
func (a *AppsUploadResource) OnStreamData(buffer []byte, b protocol.Boolean) {
    fmt.Printf("AppsUploadResource.OnStreamData - buffer %d, binary %v.\n", len(buffer), b)
}
