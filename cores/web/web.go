package web

import (
	"io"
	"net/http"
	"os"

	bz "github.com/alackfeng/bytezero/bytezero"
	"github.com/alackfeng/bytezero/cores/utils"
	"github.com/gin-gonic/gin"
)

var logweb = utils.Logger(utils.Fields{"animal": "web"})

// 允许自定义头进行跨域请求后.
const BytezeroHead = "Bytezero"
const AppHead = "App"
const AppCredentialsHead = "Credential"
const customeHeads = ", " + BytezeroHead + ", " + AppHead + ", " + AppCredentialsHead

// MIMEPROTOBUF - protobuf支持.
const MIMEPROTOBUF = "application/x-protobuf"


// GinWeb - Web Gin Gateway.GinWeb
type GinWeb struct {
    host string
    heart int32
    *gin.Engine
    bzn bz.BZNet
}

// NewGinWeb -
func NewGinWeb(uri string, ht int32, bzn bz.BZNet) *GinWeb {
    gw := &GinWeb{host: uri, heart: ht, bzn: bzn}
    // gw.Engine = gin.Default()
    gw.Engine = gin.New()
    return gw.Middlewares()
}

// Cors -
func (gw *GinWeb) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("origin")
		if len(origin) == 0 {
			origin = c.Request.Header.Get("Origin")
		}
		logweb.Debugln("Cors handle.", origin)
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"+customeHeads)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
		// c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		// c.Writer.Header().Set("Content-Type", MIMEPROTOBUF)
		// c.Writer.Header().Set("Content-Type", "application/octet-stream;charset=utf-8")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// Middlewares -
func (gw *GinWeb) Middlewares() *GinWeb {
    gw.Use(gin.Logger())
    gw.Use(gw.Cors())
	gw.Use(gin.Recovery())
	// add auth validate.
	// gw.Use(gin.BasicAuth(gin.Accounts{"username": "feng", "password": "yue"}))

    return gw
}

// Logger -
func (gw *GinWeb) Logger(logfile string) {
	if logfile == "" {
		return
	}

	f, _ := os.Create(logfile)
	// gin.DisableConsoleColor()
	// gin.DefaultWriter = io.MultiWriter(f)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

// Admins -
func (gw *GinWeb) Admins() {
}

// Statics - Statics
func (gw *GinWeb) Statics() {

	// 上传文件大小.
	gw.MaxMultipartMemory = 8 << 20 // 8M.

    // http://192.168.90.162:7790/public/upload.html
    // http://192.168.90.162:7790/public/index.html
	gw.LoadHTMLGlob("./public/*.html")
	gw.Static("/static", "./public/static")
	gw.Static("/public", "./public/")

}

// Routers - Routers
func (gw *GinWeb) Routers() {
	gw.Any("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "hello Go"})
	})
	gw.Statics()

	gw.Admins()

	v1 := gw.Group("/api/v1")
	{
		gw.RouterBridge(v1)
	}
	// gw.RouterWS()
}

// Start - Start
func (gw *GinWeb) Start() {

	gw.Routers()

	if gw.host == "" {
		logweb.Panic("GinWeb uri host is null.")
	}
	logweb.Println("GinWeb uri host is", gw.host, " , Start Now.")
	host := gw.host
	gw.Run(host)
}


