package web

import (
	"net/http"

	"github.com/alackfeng/bytezero/cores/utils"
	"github.com/gin-gonic/gin"
)



const (
	urlBaseBridge                = "/bridge"
	urlBridgeCredentialGet       = urlBaseBridge + "/credential/get"
)


// RouterBridge -
func (gw *GinWeb) RouterBridge(grg *gin.RouterGroup) {
	grg.Any(urlBridgeCredentialGet, gw.HandleBridgeCredentialGet)
}


// CredentialGetReq -
type CredentialGetReq struct {
}

// CredentialResult -
type CredentialResult struct {
    User        string    `form:"User" json:"User" xml:"User" bson:"User" binding:"required"`
    Pass        string    `form:"Pass" json:"Pass" xml:"Pass" bson:"Pass" binding:"required"`
    Expired     int64    `form:"Expired" json:"Expired" xml:"Expired" bson:"Expired" binding:"required"`
}

// HandleBridgeCredentialGet - 获取认证信息.
// http://192.168.90.162:7790/api/v1/bridge/credential/get
func (gw *GinWeb) HandleBridgeCredentialGet(c *gin.Context) {
	// gw.HandleAction(Module_api, Operator{}, c)
    result := &CredentialResult{
        Expired: utils.NowMs() + gw.bzn.CredentialExpiredMs(),
    }
    cred := utils.NewCredential(result.Expired)
    result.User = cred.Username()
    result.Pass = cred.Sign(gw.bzn.AppKey())
    c.JSON(http.StatusOK, result)
}


