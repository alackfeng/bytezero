package web

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	bzweb "github.com/alackfeng/bytezero/bytezero/web"
	"github.com/alackfeng/bytezero/cores/utils"
	"github.com/gin-gonic/gin"
)



const (
	urlBaseBridge                = "/bridge"
	urlBridgeCredentialGet       = urlBaseBridge + "/credential/get"
    urlBridgeFilesUpload         = urlBaseBridge + "/files/upload"
)


// RouterBridge -
func (gw *GinWeb) RouterBridge(grg *gin.RouterGroup) {
	grg.Any(urlBridgeCredentialGet, gw.HandleBridgeCredentialGet)
	grg.POST(urlBridgeFilesUpload, gw.HandleBridgeFilesUpload)
}


// HandleBridgeCredentialGet - 获取认证信息.
// http://192.168.90.162:7790/api/v1/bridge/credential/get
func (gw *GinWeb) HandleBridgeCredentialGet(c *gin.Context) {
	// gw.HandleAction(Module_api, Operator{}, c)

    result := bzweb.CredentialUrlResult{}
    now := utils.NowMs() + gw.bzn.CredentialExpiredMs()
    urls := gw.bzn.CredentialUrls()
    for _, url := range urls {
        credential := bzweb.CredentialURL{
            URL: url,
            Expired: now,
        }
        cred := utils.NewCredential(credential.Expired)
        credential.User = cred.Username()
        credential.Pass = cred.Sign(gw.bzn.AppKey())
        result = append(result, credential)
    }
    c.JSON(http.StatusOK, result)
}

type FileUpload struct {
    FilePath string `form:"path" json:"path" xml:"path" bson:"path" binding:"required"`
    FileName string `form:"name" json:"name" xml:"name" bson:"name" binding:"required"`
}

// SaveFileName -
func (a *FileUpload) SaveFileName() string {
    return filepath.Join(a.FilePath, a.FileName)
}

// HandleBridgeFilesUpload - 普通文件上传.
// http://192.168.90.162:7790/api/v1/bridge/files/upload
func (gw *GinWeb) HandleBridgeFilesUpload(c *gin.Context) {
    if err := CheckRequest(c, http.MethodPost, gin.MIMEMultipartPOSTForm); err != nil {
		logweb.Warnln("illegal request => ", c.Request.URL, " => ", c.Request.Method, c.ContentType())
		c.String(http.StatusNotFound, "illegal request")
    }
    form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File Form err: %s", err.Error()))
		return
	}
    upload := &FileUpload{
        FilePath: gw.uploadPath,
    }
    files := form.File["files"]
    fmt.Printf("Bytezero api upload file req: author<%s> name<%s> total<%s> files<%d> \n", c.PostForm("author"), c.PostForm("name"), c.PostForm("total"), len(files))
    for _, file := range files {

        filename := filepath.Base(file.Filename)
		if filename != "" {
			upload.FileName = filename
		} else {
            logweb.Warningln("upload error filename is nul.")
			continue
		}
		if err := utils.DirIsExistThenMkdir(upload.FilePath); err != nil {
            logweb.Warningln("upload error, mkdir savepath.", err.Error())
			continue
		}
        saveFileName := upload.SaveFileName()
        if ok, _ := utils.IsExistFile(saveFileName); ok {
            logweb.Warningln("upload error, file exists remove it now.")
            os.Remove(saveFileName)
        }
        if err := c.SaveUploadedFile(file, saveFileName); err != nil {
            logweb.Warningln("upload error, save file quit.", err.Error())
            continue
        }
    }
    c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully %d files with save to path<%s>.", len(files), upload.FilePath))
}


