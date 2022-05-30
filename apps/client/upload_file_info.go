package client

import (
	"encoding/json"
	"fmt"
)

// UploadFileInfo -
type UploadFileInfo struct {
    FilePath string `form:"FilePath" json:"FilePath" xml:"FilePath" bson:"FilePath" binding:"required"`
    FileName string `form:"FileName" json:"FileName" xml:"FileName" bson:"FileName" binding:"required"`
    FileSize int64 `form:"FileSize" json:"FileSize" xml:"FileSize" bson:"FileSize" binding:"required"`
    FileMd5 string `form:"FileMd5" json:"FileMd5" xml:"FileMd5" bson:"FileMd5" binding:"required"`
}

// NewUploadFileInfo -
func NewUploadFileInfo() *UploadFileInfo {
    return &UploadFileInfo{}
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
