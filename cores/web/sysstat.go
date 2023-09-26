package web

import (
	"net/http"

	bz "github.com/alackfeng/bytezero/bytezero"
	bzweb "github.com/alackfeng/bytezero/bytezero/web"
	"github.com/gin-gonic/gin"
)

const (
	urlBaseSysStat         = "/sysstat"
	urlSysStatNetSubscribe = urlBaseSysStat + "/net/subscribe"
	urlSysStatNetList      = urlBaseSysStat + "/net/list"
)

// RouterBridge -
func (gw *GinWeb) RouterSysStat(grg *gin.RouterGroup) {
	grg.POST(urlSysStatNetSubscribe, gw.HandleSysStatNetSubscribe)
	grg.GET(urlSysStatNetList, gw.HandleSysStatNetList)
}

// HandleSysStatNetSubscribe - 获取系统网络带宽速率.
// http://192.168.90.162:7790/api/v1/sysstat/net/subscribe
func (gw *GinWeb) HandleSysStatNetSubscribe(c *gin.Context) {
	result := bzweb.NewActionResult()
	data, err := gw.bzn.Stats(bz.SysStatNetSubscribe)
	if err != nil {
		result.Set(bzweb.ErrCode_error, err.Error())
	} else {
		result.SetData(data)
	}

	c.JSON(http.StatusOK, result)
}

// HandleSysStatNetList - 获取网卡列表.
// http://192.168.90.162:7790/api/v1/sysstat/net/list
func (gw *GinWeb) HandleSysStatNetList(c *gin.Context) {
	result := bzweb.NewActionResult()
	data, err := gw.bzn.Stats(bz.SysStatNetList)
	if err != nil {
		result.Set(bzweb.ErrCode_error, err.Error())
	} else {
		result.SetData(data)
	}

	c.JSON(http.StatusOK, result)
}
