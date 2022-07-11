package web

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	bz "github.com/alackfeng/bytezero/bytezero"
	"github.com/alackfeng/bytezero/cores/utils"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"github.com/unrolled/secure"
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

    http bool
    // https api.
    https bool
    address string
    certFile string
    keyFile string

    // files.
    uploadPath string
    logPath string
    memory int64 // = 8 << 20 // 8M.
}

// NewGinWeb -
func NewGinWeb(uri string, ht int32, bzn bz.BZNet) *GinWeb {
    gw := &GinWeb{host: uri, heart: ht, bzn: bzn }
    gw.http = uri != ""
    gw.memory = 8 << 20
    // gw.Engine = gin.Default()
    gw.Engine = gin.New()
    return gw
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

// AccessIps -
func (gw *GinWeb) AccessIps() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.RemoteIP()
        if err := gw.bzn.AccessIpsAllow(ip); err != nil {
            logweb.Errorln("Web.AccessIpsAllow - ip ", ip, err.Error())
            c.AbortWithStatus(200)
            return
        }
        if err := gw.bzn.AccessIpsDeny(ip); err != nil {
            logweb.Errorln("Web.AccessIpsDeny - ip ", ip, err.Error())
            c.AbortWithStatus(200)
            return
        }
        if err := gw.bzn.AccessIpsForbid(ip, false); err != nil {
            logweb.Errorln("Web.AccessIpsForbid - ip ", ip, err.Error())
        }
        c.Next()
    }
}

// Middlewares -
func (gw *GinWeb) Middlewares() *GinWeb {
    gw.Logger(gw.logPath)
    gw.Use(gin.Logger())

    gw.Use(gw.Cors())
	gw.Use(gin.Recovery())
    gw.Use(gw.AccessIps())
	// add auth validate.
	// gw.Use(gin.BasicAuth(gin.Accounts{"username": "feng", "password": "yue"}))
    return gw
}

// Logger -
func (gw *GinWeb) Logger(logPath string) {
    logfile := filepath.Join(gw.logPath, utils.LogName("bytezero_access"))
    logweb.Infoln("Web.Logger to ", logfile)
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
	gw.MaxMultipartMemory = gw.memory

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

    gw.Use(favicon.New("./public/favicon.ico"))



	// gw.RouterWS()
}

// TlsHandler -
func (gw *GinWeb) TlsHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        secureMiddleware := secure.New(secure.Options{
            SSLRedirect: true,
            SSLHost:     "localhost:8080",
        })
        err := secureMiddleware.Process(c.Writer, c.Request)
        if err != nil {
            logweb.Errorln("Web.TlsHandler - error", err.Error())
            return
        }
        logweb.Errorln("Web.TlsHandler - ", c.RemoteIP())
        c.Next()
    }
}

// SetSecretTransport -
func (gw *GinWeb) SetSecretTransport(address string, certFile string, keyFile string) *GinWeb {
    gw.address = address
    gw.certFile = certFile
    gw.keyFile = keyFile
    if gw.address == "" || gw.certFile == "" || gw.keyFile == "" {
        logweb.Errorf("Web.SetSecretTransport - must is null. maybe %s, %s, %s.", gw.address, gw.certFile, gw.keyFile)
        return gw
    }
    gw.https = true
    return gw
}

// SetStaticInfo -
func (gw *GinWeb) SetStaticInfo(memory int64, logPath string, uploadPath string) *GinWeb {
    gw.logPath = logPath
    gw.uploadPath = uploadPath
    gw.memory = memory
    utils.DirIsExistThenMkdir(filepath.Join(gw.logPath))
    utils.DirIsExistThenMkdir(filepath.Join(gw.uploadPath))
    return gw
}

// StartTls -
func (gw *GinWeb) StartTls() {
    if !gw.https {
        return
    }
    if gw.address == "" || gw.certFile == "" || gw.keyFile == "" {
        logweb.Errorf("Web.StartTls - must is null. maybe %s, %s, %s.", gw.address, gw.certFile, gw.keyFile)
        return
    }

    gw.Use(gw.TlsHandler())
    logweb.Println("Web.Start https://", gw.address, " , Start Now.")
    gw.RunTLS(gw.address, gw.certFile, gw.keyFile)
    logweb.Println("Web.Start https://", gw.address, " , Start Over.")
}

// Start - Start
func (gw *GinWeb) Start() {
    gw.Middlewares()
	gw.Routers()

    if gw.https && gw.http {
        go gw.StartTls()
    } else if gw.https {
        gw.StartTls()
    }

    if gw.http {
        if gw.host == "" {
            logweb.Panic("Web.Start uri host is null.")
        }
        logweb.Println("Web.Start uri host is", gw.host, " , Start Now.")
        host := gw.host
        gw.Run(host)
        logweb.Println("Web.Start uri host is", gw.host, " , Start Over.")
    }
}


