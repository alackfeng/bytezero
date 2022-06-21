package client

import (
	"encoding/json"
	"fmt"

	"github.com/alackfeng/bytezero/cores/utils"
)



const (
    CmdNameCreateUploadTask = "createUploadTask"
    CmdNameFinishUploadTask = "finishUploadTask"
    CmdNameReceivedFinishMessage = "receivedFinishMessage"
)

// UploadFileInfo -
type UploadFileInfo struct {
    CmdName string `form:"cmdName" json:"cmdName" xml:"cmdName" bson:"cmdName" binding:"required"`
    TaskId string `form:"taskId" json:"taskId" xml:"taskId" bson:"taskId" binding:"required"`
    FilePath string `form:"path" json:"path" xml:"path" bson:"path" binding:"required"`
    FileName string `form:"name" json:"name" xml:"name" bson:"name" binding:"required"`
    FileSize int64 `form:"len" json:"len" xml:"len" bson:"len" binding:"required"`
    FileMd5 string `form:"md5" json:"md5" xml:"md5" bson:"md5" binding:"required"`
}

// NewUploadFileInfo -
func NewUploadFileInfo(cmdName string) *UploadFileInfo {
    return &UploadFileInfo{
        CmdName: cmdName,
        TaskId: utils.UUID(),
    }
}

// To -
func (u *UploadFileInfo) To() []byte {
    mByte, _ := json.Marshal(u)
    return mByte
}

func (u *UploadFileInfo) ToMd5() []byte {
    mByte, _ := json.Marshal(u)
    return mByte
}

// From -
func (u *UploadFileInfo) From(b []byte) error {
    return json.Unmarshal(b, u)
}

// String -
func (u UploadFileInfo) String() string {
    return fmt.Sprintf("FileInfo[%s,%s,%d]", u.FilePath, u.FileName, u.FileSize)
}
