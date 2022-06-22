package client

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/alackfeng/bytezero/bytezero/protocol"
	"github.com/alackfeng/bytezero/cores/utils"
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

    f *os.File
    info UploadFileInfo
    f5 hash.Hash
    offset int

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
func NewAppsUploadResourceAnswer(app Client, sessionId string, savePath string) *AppsUploadResource {
    return &AppsUploadResource{
        app: app,
        sessionId: sessionId,
        filePath: savePath,
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

    f, err := os.Open(a.filePath)
    if err != nil {
        fmt.Printf("AppsUploadResource.uploadFile - upload file<%s>, error.%s\n", a.filePath, err.Error())
        return err
    }
    defer f.Close()

    fstat, _ := f.Stat()
    a.info = UploadFileInfo{
        CmdName: CmdNameCreateUploadTask,
        TaskId: utils.UUID(),
        FilePath: a.filePath,
        FileName: fstat.Name(),
        FileSize: fstat.Size(),
    }

    // 创建Stream通道.
    a.streamHandle, err = a.channelHandle.StreamCreate(a, a.info.To())
    if err != nil {
        fmt.Printf("AppsUploadResource.uploadFile - streamHandle is null, error.%v.\n", err.Error())
        return err
    }

    fmt.Printf("AppsUploadResource.uploadFile - begin upload file<%s>, at Channel#%dStream#%d.\n", a.filePath, a.channelHandle.Id(), a.streamHandle.StreamId())

    offset := 0
    a.f5 = md5.New()
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
        if (offset / a.bufferLen) % 1000 == 0 {
            fmt.Printf("AppsUploadResource.uploadFile - send to Channel#%dStream#%d, buffer %d, offset %d - %d.\n", a.channelHandle.Id(), a.streamHandle.StreamId(), n, offset, a.info.FileSize)
        }
        if err := a.streamHandle.SendData(buf[0:n]); err != nil {
            return err
        }
        offset += n
        a.f5.Write(buf[0:n])
        time.Sleep(time.Millisecond * 5)
    }
    a.info.FileMd5 = fmt.Sprintf("%X", a.f5.Sum(nil))
    fmt.Printf("AppsUploadResource.uploadFile - end.. upload file<%s> size<%d> md5<%s>, at Channel#%dStream#%d over.\n", a.filePath, a.info.FileSize, a.info.FileMd5, a.channelHandle.Id(), a.streamHandle.StreamId())
    a.info.CmdName = CmdNameFinishUploadTask
    a.streamHandle.SendSignal(a.info.To())
    // a.channelHandle.StreamClose(a.streamId, a.info.ToMd5())
    // a.app.ChannelCloseByHandle(a.channelHandle)
    // go a.uploadFile()
    return nil
}

// SavePath -
func (a *AppsUploadResource) SavePath() string {
    return filepath.Join(a.filePath, fmt.Sprintf("%d_%s", a.streamHandle.StreamId(), a.info.FileName))
}

// answerFile -
func (a *AppsUploadResource) answerFile() error {
    if a.filePath == "" {
        fmt.Printf("AppsUploadResource::answerFile - save path<%s> is null.", a.filePath)
        return fmt.Errorf("save path is null")
    }
    if a.info.FileName == "" {
        fmt.Printf("AppsUploadResource::answerFile - filename<%s> is null.", a.info.FileName)
        return fmt.Errorf("filename is null")
    }
    fileName := a.SavePath()
    f, err := os.Create(fileName)
    if err != nil {
        fmt.Printf("AppsUploadResource::answerFile - Open<%s> error.%v\n", a.filePath, err.Error())

        return err
    }
    a.f = f
    a.f5 = md5.New()
    fmt.Printf("AppsUploadResource::answerFile - save to<%s>.\n", fileName)
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
    if a.mode == ModeUpload {
        go a.uploadFile()
    }
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

// OnState - ChannelObserver interface.
func (a *AppsUploadResource) OnState(state protocol.ChannelState) {
    fmt.Printf("AppsUploadResource.OnState - state %v.\n", state)
}


// onStream - ChannelObserver interface.
func (a *AppsUploadResource) onStream(s StreamHandle) (protocol.ErrCode, error) {
    fmt.Printf("AppsUploadResource.onStream - Stream#%d\n", s.StreamId())
    a.streamHandle = s
    a.streamHandle.RegisterObserver(a)
    a.info.From(a.streamHandle.ExtraInfo())
    fmt.Printf("AppsUploadResource.onStream - Stream#%d, Extra: %v.\n", s.StreamId(), a.info)
    a.offset = 0
    a.f5 = md5.New()
    return protocol.ErrCode_streamCreateAck, fmt.Errorf("ok")
}


// OnStreamSuccess -
func (a *AppsUploadResource) OnStreamSuccess(sid protocol.StreamId) {
    fmt.Println("AppsUploadResource.OnStreamSuccess - stream id ", sid)
    a.streamId = sid

}

// OnStreamError -
func (a *AppsUploadResource) OnStreamError(code int, message string) {
    fmt.Printf("AppsUploadResource.OnStreamError - code %d, %s.\n", code, message)
}

// OnStreamData -
func (a *AppsUploadResource) OnStreamData(buffer []byte, b protocol.Boolean) {
    // fmt.Printf("AppsUploadResource.OnStreamData - buffer %d, binary %v.\n", len(buffer), b)
    if b == protocol.BooleanTrue && a.mode == ModeAnswer {
        if a.f == nil {
            if err := a.answerFile(); err != nil {
                return
            }
        }
        a.offset += len(buffer)
        a.f5.Write(buffer)
        if n, err := a.f.Write(buffer); err != nil || n != len(buffer) {
            if err != nil {
                fmt.Printf("AppsUploadResource.OnStreamData - write to<%s> failed, error.%v.\n", a.f.Name(), err.Error())
            } else {
                fmt.Printf("AppsUploadResource.OnStreamData - write to<%s> failed not all writen.%d.\n", a.f.Name(), n)
            }
            return
        }

    } else {
        fmt.Printf("AppsUploadResource.OnStreamData - buffer %d, binary %v.\n", len(buffer), b)
        if a.mode == ModeAnswer {
            info := UploadFileInfo{}
            info.From(buffer)
            a.info.FileMd5 = fmt.Sprintf("%X", a.f5.Sum(nil))
            fmt.Printf("AppsUploadResource.OnStreamData - %v(%s) => total length %d:%d.\n", a.info.FilePath, info.FileMd5, a.offset, info.FileSize)
            fmt.Printf("AppsUploadResource.OnStreamData - %v(%s) =>%v.\n", a.SavePath(), a.info.FileMd5,  a.info.FileMd5 == info.FileMd5)
        } else {
            a.channelHandle.StreamClose(a.streamId, nil)
        }

    }
}

// OnStreamClosing -
func (a *AppsUploadResource) OnStreamClosing(extra []byte) {
    if a.f == nil {
        return
    }
    a.f.Close()
}
