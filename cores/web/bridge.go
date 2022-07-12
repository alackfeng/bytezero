package web

import (
	"fmt"
	"net/http"
	"net/url"
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
    urlBridgeFilesDownload       = urlBaseBridge + "/files/download"
    urlBridgeFilesList           = urlBaseBridge + "/files/list"
    urlBridgeLogsDownload        = urlBaseBridge + "/logs/download"
    urlBridgeLogsList            = urlBaseBridge + "/logs/list"
    urlBridgeAccessIpsForbid     = urlBaseBridge + "/accessips/forbid"
    urlBridgeAccessIpsReload     = urlBaseBridge + "/accessips/reload"
)


// RouterBridge -
func (gw *GinWeb) RouterBridge(grg *gin.RouterGroup) {
	grg.GET(urlBridgeCredentialGet, gw.HandleBridgeCredentialGet)
    grg.POST(urlBridgeCredentialGet, gw.HandleBridgeCredentialGet)
	grg.POST(urlBridgeFilesUpload, gw.HandleBridgeFilesUpload)
	grg.GET(urlBridgeFilesDownload, gw.HandleBridgeFilesDownload)
	grg.GET(urlBridgeFilesList, gw.HandleBridgeFilesList)
	grg.GET(urlBridgeLogsDownload, gw.HandleBridgeLogsDownload)
	grg.GET(urlBridgeLogsList, gw.HandleBridgeLogsList)
	grg.GET(urlBridgeAccessIpsForbid, gw.HandleBridgeAccessIpsForbid)
    grg.GET(urlBridgeAccessIpsReload, gw.HandleBridgeAccessIpsReload)
}


// HandleBridgeCredentialGet - 获取认证信息.
// http://192.168.90.162:7790/api/v1/bridge/credential/get
func (gw *GinWeb) HandleBridgeCredentialGet(c *gin.Context) {
	// gw.HandleAction(Module_api, Operator{}, c)

    result := bzweb.CredentialUrlResult{}
    now := utils.NowMs() + gw.bzn.CredentialExpiredMs()
    urls := gw.bzn.CredentialUrls()
    for _, u := range urls {
        ur, err := url.Parse(u)
        if err != nil {
            continue
        }
        credential := bzweb.CredentialURL{
            Scheme: ur.Scheme,
            IP: ur.Hostname(),
            Port: ur.Port(),
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

// DownFileName -
func (a *FileUpload) DownFileName() string {
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

// HandleBridgeFilesDownload -
func (gw *GinWeb) HandleBridgeFilesDownload(c *gin.Context) {

    upload := &FileUpload{
        FilePath: gw.uploadPath,
        FileName: c.Query("file"),
    }
    if err := utils.DirIsExistThenMkdir(upload.FilePath); err != nil {
        logweb.Warningln("download error, path not existed.", err.Error())
        c.String(http.StatusNotFound, fmt.Sprintf("file<%s> not found", upload.FileName))
        return
    }
    downFileName := upload.DownFileName()
    if ok, _ := utils.IsExistFile(downFileName); ok {
        logweb.Warningln("download error, file not existed.")
        c.String(http.StatusNotFound, fmt.Sprintf("file<%s> not found", upload.FileName))
        return
    }
    c.Header("Content-Type", "application/octet-stream;text/html")
    c.Header("Content-Disposition", "attachment; filename="+upload.FileName)
    c.File(downFileName)

    c.String(http.StatusOK, fmt.Sprintf("Download successfully file<%s>.", downFileName))
}


// HandleBridgeFilesList -
func (gw *GinWeb) HandleBridgeFilesList(c *gin.Context) {

}

// HandleBridgeLogsDownload -
func (gw *GinWeb) HandleBridgeLogsDownload(c *gin.Context) {
    logId := c.Query("log")
    if err := utils.DirIsExistThenMkdir(gw.logPath); err != nil {
        logweb.Warningln("log error, path not existed.", err.Error())
        c.String(http.StatusNotFound, fmt.Sprintf("log<%s> not found", logId))
        return
    }
    logFileName := filepath.Join(gw.logPath, logId)
    if ok, _ := utils.IsExistFile(logFileName); ok {
        logweb.Warningln("log error, file not existed.")
        c.String(http.StatusNotFound, fmt.Sprintf("log<%s> not found", logId))
        return
    }
    c.Header("Content-Type", "application/octet-stream;text/html")
    c.Header("Content-Disposition", "attachment; filename="+logFileName)
    c.File(logFileName)

    c.String(http.StatusOK, fmt.Sprintf("Download successfully log file<%s>.", logId))
}

// HandleBridgeLogsList -
func (gw *GinWeb) HandleBridgeLogsList(c *gin.Context) {
}


// HandleBridgeAccessIpsForbid -
// http://192.168.90.162:7790/api/v1/bridge/accessips/forbid?ip=192.168.90.162&forbid=1|2
func (gw *GinWeb) HandleBridgeAccessIpsForbid(c *gin.Context) {
    result := bzweb.NewActionResult()
    accessIpsForbidAction := &bzweb.AccessIpsForbidAction{}
    if c.Request.Method == http.MethodPost {
        if err := c.ShouldBindJSON(accessIpsForbidAction); err != nil {
            result.Set(bzweb.ErrCode_error, err.Error())
        }
    } else if c.Request.Method == http.MethodGet {
        if err := c.ShouldBindQuery(accessIpsForbidAction); err != nil {
            result.Set(bzweb.ErrCode_error, err.Error())
        }
    }

    if !accessIpsForbidAction.Check() {
        result.Set(bzweb.ErrCode_error, "nil")
    } else if err := gw.bzn.AccessIpsForbid(accessIpsForbidAction.IP, accessIpsForbidAction.Deny == 1); err != nil {
        result.Set(bzweb.ErrCode_success, err.Error())
    }
    c.JSONP(http.StatusOK, result)
}


// HandleBridgeAccessIpsForbid -
// http://192.168.90.162:7790/api/v1/bridge/accessips/reload
func (gw *GinWeb) HandleBridgeAccessIpsReload(c *gin.Context) {
    result := bzweb.NewActionResult()
    accessIpsReloadAction := &bzweb.AccessIpsReloadAction{}
    if c.Request.Method == http.MethodPost {
        if err := c.ShouldBindJSON(accessIpsReloadAction); err != nil {
            result.Set(bzweb.ErrCode_error, err.Error())
        }
    } else if c.Request.Method == http.MethodGet {
        if err := c.ShouldBindQuery(accessIpsReloadAction); err != nil {
            result.Set(bzweb.ErrCode_error, err.Error())
        }
    }

    if !accessIpsReloadAction.Check() {
        result.Set(bzweb.ErrCode_error, "nil")
    } else if err := gw.bzn.AccessIpsReload(accessIpsReloadAction.Allow()); err != nil {
        result.Set(bzweb.ErrCode_success, err.Error())
    }
    c.JSONP(http.StatusOK, result)
}
