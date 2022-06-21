package web

import (
	"net/http"

	bzweb "github.com/alackfeng/bytezero/bytezero/web"
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


// HandleBridgeCredentialGet - 获取认证信息.
// http://192.168.90.162:7790/api/v1/bridge/credential/get
func (gw *GinWeb) HandleBridgeCredentialGet(c *gin.Context) {
	// gw.HandleAction(Module_api, Operator{}, c)

    result := &bzweb.CredentialResult{
        Expired: utils.NowMs() + gw.bzn.CredentialExpiredMs(),
    }
    cred := utils.NewCredential(result.Expired)
    result.User = cred.Username()
    result.Pass = cred.Sign(gw.bzn.AppKey())
    c.JSON(http.StatusOK, result)
}


