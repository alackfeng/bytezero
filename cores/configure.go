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
    IP string `yaml:"ip" json:"ip" binding:"required"`
    Port int `yaml:"port" json:"port" binding:"required"`
}

// Address -
func (c *AppServerConfigure) Address() string {
    if c.IP == "" {
        return fmt.Sprintf(":%d", c.Port)
    }
    return fmt.Sprintf("%s:%d", c.IP, c.Port)
}

// AppWebConfigure -
type AppWebConfigure struct {
    Host string `yaml:"host" json:"host" binding:"required"`
    Heart int32 `yaml:"heart" json:"heart" binding:"required"`
}

// AppCredentialConfig -
type AppCredentialConfig struct {
    ExpiredMs int64 `yaml:"expiredMs" json:"expiredMs" binding:"required"`
}

// AppConfigure -
type AppConfigure struct {
    Name string  `yaml:"name" json:"name" binding:"required"`
    Version string `yaml:"version" json:"version" binding:"required"`
    Server AppServerConfigure `yaml:"server" json:"server" binding:"required"`
    Web AppWebConfigure `yaml:"web" json:"web" binding:"required"`
    MaxBufferLen int `yaml:"maxBufferLen" json:"maxBufferLen" binding:"required"`
    RWBufferLen int `yaml:"rwBufferLen" json:"rwBufferLen" binding:"required"`
    Appid string `yaml:"appid" json:"appid" binding:"required"`
    Appkey string `yaml:"appkey" json:"appkey" binding:"required"`
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
    return fmt.Sprintf("\n App\t: %s-%s, ID: %s, Key:%s, MaxBuffer: %d, RWBuffer: %d \n Server\t: tcp://%s \n Web\t: http://%s(heart:%d) \n",
        c.App.Name, c.App.Version, c.App.Appid, c.App.Appkey, c.App.MaxBufferLen, c.App.RWBufferLen,
        c.App.Server.Address(),
        c.App.Web.Host, c.App.Web.Heart)
}

// ConfigSetServer -
func ConfigSetServer(maxBufferLen, rwBufferLen, port int, host, appid, appkey string) {
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
        GlobalConfig.App.Web.Host = host
    }
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
