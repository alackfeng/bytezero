package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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

// CheckRequest -
func CheckRequest(c *gin.Context, method, mimetype string) error {
	if c.Request.Method != method {
		return fmt.Errorf("Method %s No Support", c.Request.Method)
	}
	if c.ContentType() != mimetype {
		return fmt.Errorf("Content-Type %s No Support", c.ContentType())
	}
	return nil
}

// HandleAction -
func (gw *GinWeb) HandleAction(m Module, oper Operator, c *gin.Context) {

}
