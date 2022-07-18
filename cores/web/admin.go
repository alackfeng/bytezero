package web

import (
	"net/http"

	bzweb "github.com/alackfeng/bytezero/bytezero/web"
	"github.com/gin-gonic/gin"
)

const (
	urlBaseAdmin                = "/admin"
    urlAdminAccessIpsForbid     = urlBaseAdmin + "/accessips/forbid"
    urlAdminAccessIpsReload     = urlBaseAdmin + "/accessips/reload"
    urlAdminSystemRestart       = urlBaseAdmin + "/system/restart"
    urlAdminSystemStop          = urlBaseAdmin + "/system/stop"
    urlAdminSystemReload        = urlBaseAdmin + "/system/reload"
)

// RouterAdmin -
func (gw *GinWeb) RouterAdmin(grg *gin.RouterGroup) {
	grg.GET(urlAdminAccessIpsForbid, gw.HandleAdminAccessIpsForbid)
    grg.GET(urlAdminAccessIpsReload, gw.HandleAdminAccessIpsReload)
    grg.GET(urlAdminSystemRestart, gw.HandleAdminSystemRestart)
    grg.GET(urlAdminSystemStop, gw.HandleAdminSystemStop)
    grg.GET(urlAdminSystemReload, gw.HandleAdminSystemReload)
}


// HandleBridgeAccessIpsForbid -
// http://192.168.90.162:7790/api/v1/admin/accessips/forbid?ip=192.168.90.162&forbid=1|2
func (gw *GinWeb) HandleAdminAccessIpsForbid(c *gin.Context) {
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
// http://192.168.90.162:7790/api/v1/admin/accessips/reload
func (gw *GinWeb) HandleAdminAccessIpsReload(c *gin.Context) {
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

// HandleAdminSystemRestart -
func (gw *GinWeb) HandleAdminSystemRestart(c *gin.Context) {
    result := bzweb.NewActionResult()
    if err := gw.bzn.SystemRestart(); err != nil {
        result.Set(bzweb.ErrCode_error, err.Error())
    }
    c.JSONP(http.StatusOK, result)
}

// HandleAdminSystemStop -
func (gw *GinWeb) HandleAdminSystemStop(c *gin.Context) {
    result := bzweb.NewActionResult()
    if err := gw.bzn.SystemStop(); err != nil {
        result.Set(bzweb.ErrCode_error, err.Error())
    }
    c.JSONP(http.StatusOK, result)
}

// HandleAdminSystemReload -
func (gw *GinWeb) HandleAdminSystemReload(c *gin.Context) {
    result := bzweb.NewActionResult()
    if err := gw.bzn.SystemReload(); err != nil {
        result.Set(bzweb.ErrCode_error, err.Error())
    }
    c.JSONP(http.StatusOK, result)
}
