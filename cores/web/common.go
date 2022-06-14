package web

import "github.com/gin-gonic/gin"

// Module -
type Module int32

const (
	Module_none    Module = 0
	Module_api     Module = 1
)

// Operator -
type Operator struct {
    IP    string   // 公网IP地址.
    AppID string
    Token string   // AccessToken or RefreshToken.
}

// HandleAction -
func (gw *GinWeb) HandleAction(m Module, oper Operator, c *gin.Context) {

}
