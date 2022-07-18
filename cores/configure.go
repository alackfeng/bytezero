package cores

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// GlobalConfig -
var GlobalConfig Configure

// ConfigGlobal - 全局实例
func ConfigGlobal() *Configure {
	return &GlobalConfig
}

// APPVersion APP本地控制
var APPVersion = "v1.0.0"

// AppServerConfigure -
type AppServerConfigure struct {
    UP bool `yaml:"up" json:"up" binding:"required"`
    IP string `yaml:"ip" json:"ip" binding:"required"`
    Port int `yaml:"port" json:"port" binding:"required"`
    Margic bool `yaml:"margic" json:"margic" binding:"required"`
}

// Address -
func (c *AppServerConfigure) Address() string {
    if c.IP == "" {
        return fmt.Sprintf(":%d", c.Port)
    }
    return fmt.Sprintf("%s:%d", c.IP, c.Port)
}

// AppTlsConfigure -
type AppTlsConfigure struct {
    UP bool `yaml:"up" json:"up" binding:"required"`
    IP string `yaml:"ip" json:"ip" binding:"required"`
    Port int `yaml:"port" json:"port" binding:"required"`
    CaCert string `yaml:"cacert" json:"cacert" binding:"required"`
    CaKey string `yaml:"cakey" json:"cakey" binding:"required"`
}

// Address -
func (c *AppTlsConfigure) Address() string {
    if c.IP == "" {
        return fmt.Sprintf(":%d", c.Port)
    }
    return fmt.Sprintf("%s:%d", c.IP, c.Port)
}

// AppWebHttpConfigure -
type AppWebHttpConfigure struct {
    UP bool `yaml:"up" json:"up" binding:"required"`
    Host string `yaml:"host" json:"host" binding:"required"`
    Heart int32 `yaml:"heart" json:"heart" binding:"required"`
}

// Address -
func (a AppWebHttpConfigure) Address() string {
    return fmt.Sprintf("http://%s(heart:%d, up:%v)", a.Host, a.Heart, a.UP)
}

// AppWebHttpsConfigure -
type AppWebHttpsConfigure struct {
    UP bool `yaml:"up" json:"up" binding:"required"`
    Host string `yaml:"host" json:"host" binding:"required"`
    Heart int32 `yaml:"heart" json:"heart" binding:"required"`
    CaCert string `yaml:"cacert" json:"cacert" binding:"required"`
    CaKey string `yaml:"cakey" json:"cakey" binding:"required"`
}

// Address -
func (a AppWebHttpsConfigure) Address() string {
    return fmt.Sprintf("https://%s(heart:%d, up:%v)", a.Host, a.Heart, a.UP)
}

// AppWebStaticConfigure -
type AppWebStaticConfigure struct {
    UP bool `yaml:"up" json:"up" binding:"required"`
    UploadPath string `yaml:"uploadPath" json:"uploadPath" binding:"required"`
    LogPath string `yaml:"logPath" json:"logPath" binding:"required"`
    Memory int64 `yaml:"memory" json:"memory" binding:"required"`
}

// AppWebConfigure -
type AppWebConfigure struct {
    Http AppWebHttpConfigure `yaml:"http" json:"http" binding:"required"`
    Https AppWebHttpsConfigure `yaml:"https" json:"https" binding:"required"`
    Static AppWebStaticConfigure `yaml:"static" json:"static" binding:"required"`
}

// AppCredentialConfig -
type AppCredentialConfig struct {
    ExpiredMs int64 `yaml:"expiredMs" json:"expiredMs" binding:"required"`
    Urls []string `yaml:"urls" json:"urls" binding:"required"`
}

// AppConfigure -
type AppConfigure struct {
    Name string  `yaml:"name" json:"name" binding:"required"`
    Version string `yaml:"version" json:"version" binding:"required"`
    Server AppServerConfigure `yaml:"server" json:"server" binding:"required"`
    Tls AppTlsConfigure `yaml:"tls" json:"tls" binding:"required"`
    Web AppWebConfigure `yaml:"web" json:"web" binding:"required"`
    MaxBufferLen int `yaml:"maxBufferLen" json:"maxBufferLen" binding:"required"`
    RWBufferLen int `yaml:"rwBufferLen" json:"rwBufferLen" binding:"required"`
    Appid string `yaml:"appid" json:"appid" binding:"required"`
    Appkey string `yaml:"appkey" json:"appkey" binding:"required"`
    LogPath string `yaml:"logPath" json:"logPath" binding:"required"`
    AccessIpsAllow string `yaml:"accessIpsAllow" json:"accessIpsAllow" binding:"required"`
    AccessIpsDeny string `yaml:"accessIpsDeny" json:"accessIpsDeny" binding:"required"`
    Credential AppCredentialConfig `yaml:"credential" json:"credential" binding:"required"`
}

// Configure -
type Configure struct {
    App AppConfigure `yaml:"app" json:"app" binding:"required"`
}

// NewConfigure -
func NewConfigure() *Configure {
    return &Configure{}
}

// Version -
func (c* Configure) Version() string {
    if c.App.Version != "" {
        fmt.Fprintf(os.Stderr, "bytezero go APP Version: %s\n", c.App.Version)
        return c.App.Version
    }
    return APPVersion
}

// String -
func (c Configure) String() string {
    return fmt.Sprintf("\n App\t: %s-%s, ID: %s, Key:%s, MaxBuffer: %d, RWBuffer: %d \n Server\t: tcp://%s \n Web\t: %s, %s. \n",
        c.App.Name, c.App.Version, c.App.Appid, c.App.Appkey, c.App.MaxBufferLen, c.App.RWBufferLen,
        c.App.Server.Address(),
        c.App.Web.Http.Address(), c.App.Web.Https.Address())
}

// ConfigSetServer -
func ConfigSetServer(maxBufferLen, rwBufferLen, port int, host, appid, appkey string, margic bool) {
    if GlobalConfig.App.Name != "" {
        return
    }
    if maxBufferLen != 0 {
        GlobalConfig.App.MaxBufferLen = maxBufferLen
    }
    if rwBufferLen != 0 {
        GlobalConfig.App.RWBufferLen = rwBufferLen
    }
    if appid != "" {
        GlobalConfig.App.Appid = appid
    }
    if appkey != "" {
        GlobalConfig.App.Appkey = appkey
    }
    if port != 0 {
        GlobalConfig.App.Server.Port = port
    }
    if host != "" {
        GlobalConfig.App.Web.Http.Host = host
    }

    GlobalConfig.App.Server.Margic = margic
}

// ConfigSetTls -
func ConfigSetTls(needTls bool, tlsPort int, caCert string, caKey string) {
    if GlobalConfig.App.Name != "" {
        return
    }
    if !needTls {
        return
    }
    GlobalConfig.App.Tls.UP = needTls
    GlobalConfig.App.Tls.Port = tlsPort
    GlobalConfig.App.Tls.CaCert = caCert
    GlobalConfig.App.Tls.CaKey = caKey
}

// ConfigureParse
func ConfigureParse() (*Configure, error) {
    conf := &Configure{}
    if err := viper.Unmarshal(conf); err != nil {
        fmt.Println("configure load falied.", err.Error())
        return nil, nil
    } else {
        fmt.Println("configure load: ", conf)
    }
    GlobalConfig = *conf
    return conf, nil
}
