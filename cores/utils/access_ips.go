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
    return yaml.NewDecoder(f).Decode(a)
}

// Upload -
func (a *AccessIpsDeny) Upload(accessIpFile string, ip string) error {
    if accessIpFile == "" {
        fmt.Println("AccessIpsDeny::Upload accessIpFile is nil.")
        return nil
    }

    f, err := os.OpenFile(accessIpFile, os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("AccessIpsDeny::Upload f is nil.", err.Error())
        return err
    }

    a.Deny[ip] = true // add one.

    writer := bufio.NewWriter(f)
    writer.WriteString(fmt.Sprintf("  %s: true", ip))
    writer.Flush()
    f.Close()
    return nil
}



