package utils

import (
	"bufio"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// AccessIpsAllow -
type AccessIpsAllow struct {
    Allow map[string]bool `yaml:"allow,flow" json:"allow,flow" binding:"-"`
}

// Load -
func (a *AccessIpsAllow) Load(accessIpFile string) error {
    if accessIpFile == "" {
        return nil
    }
    f, err := os.Open(accessIpFile)
    if err != nil {
        return err
    }
    defer f.Close()
    return yaml.NewDecoder(f).Decode(a)
}

// AccessIpsDeny -
type AccessIpsDeny struct {
    Deny map[string]bool `yaml:"deny,flow" json:"deny,flow" binding:"-"`
}

// Load -
func (a *AccessIpsDeny) Load(accessIpFile string) error {
    if accessIpFile == "" {
        return nil
    }
    f, err := os.Open(accessIpFile)
    if err != nil {
        return err
    }
    defer f.Close()
    a.Deny = make(map[string]bool)
    return yaml.NewDecoder(f).Decode(a)
}

// Upload -
func (a *AccessIpsDeny) Upload(accessIpFile string, ip string, deny bool) error {
    if accessIpFile == "" {
        fmt.Println("AccessIpsDeny::Upload accessIpFile is nil.")
        return nil
    }

    f, err := os.OpenFile(accessIpFile, os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("AccessIpsDeny::Upload f is nil.", err.Error())
        return err
    }

    // 已经存在不更新，需要手动修改配置并重启.
    if denyIp, ok := a.Deny[ip]; ok {
        fmt.Println("----------AccessIpsDeny ok, denyIp ", ip, denyIp, deny)
        // if denyIp == deny {
            return nil
        // }
    }
    fmt.Println("----------AccessIpsDeny add ", ip, deny)
    a.Deny[ip] = deny // add one.

    writer := bufio.NewWriter(f)
    writer.WriteString(fmt.Sprintf("  %s: %v\n", ip, deny))
    writer.Flush()
    f.Close()
    return nil
}



