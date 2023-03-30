package web

import (
	"net/http"

	bzweb "github.com/alackfeng/bytezero/bytezero/web"
	"github.com/gin-gonic/gin"
)

const (
	urlBaseSysStat             = "/sysstat"
	urlBridSysStatNetSubscribe = urlBaseSysStat + "/net/subscribe"
)

// RouterBridge -
func (gw *GinWeb) RouterSysStat(grg *gin.RouterGroup) {
	grg.GET(urlBridSysStatNetSubscribe, gw.HandleBridSysStatNetSubscribe)
}

// CredentialURL -
type CredentialURL struct {
	Scheme  string `form:"Scheme" json:"Scheme" xml:"Scheme" bson:"Scheme" binding:"required"`
	IP      string `form:"IP" json:"IP" xml:"IP" bson:"IP" binding:"required"`
	Port    string `form:"Port" json:"Port" xml:"Port" bson:"Port" binding:"required"`
	User    string `form:"User" json:"User" xml:"User" bson:"User" binding:"required"`
	Pass    string `form:"Pass" json:"Pass" xml:"Pass" bson:"Pass" binding:"required"`
	Expired int64  `form:"ExpiredMs" json:"ExpiredMs" xml:"ExpiredMs" bson:"ExpiredMs" binding:"required"`
}

// HandleBridSysStatNetSubscribe - 获取系统网络带宽速率.
// http://192.168.90.162:7790/api/v1/sysstat/net/subscribe
func (gw *GinWeb) HandleBridSysStatNetSubscribe(c *gin.Context) {
	result := bzweb.NewActionResult()
	data, err := gw.bzn.Stats()
	if err != nil {
		result.Set(bzweb.ErrCode_error, err.Error())
	} else {
		result.SetData(data)
	}

	c.JSON(http.StatusOK, result)
}
