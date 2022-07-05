package server

import (
	"crypto/rand"
	"crypto/tls"
	"time"

	bz "github.com/alackfeng/bytezero/bytezero"
)

// TlsServer -
type TlsServer struct {
    bzn bz.BZNet
    network string
    address string
    certFile string
    keyFile string
}

// NewTlsServer -
func NewTlsServer(bzn bz.BZNet, address string, certFile, keyFile string) *TlsServer {
    return &TlsServer{
        bzn: bzn,
        network: "tcp",
        address: address,
        certFile: certFile,
        keyFile: keyFile,
    }
}

// Listen -
func (t *TlsServer) Listen() error {
    crt, err := tls.LoadX509KeyPair(t.certFile, t.keyFile)
    if err != nil {
        return err
    }
    config := &tls.Config{
        Certificates: []tls.Certificate{crt},
        Time: time.Now,
        Rand: rand.Reader,
    }
    listen, err := tls.Listen(t.network, t.address, config)
    if err != nil {
        return err
    }
    logsv.Debugln("====>TlsServer Listen TLS:", t.address)

    for {
        conn, err := listen.Accept()
        if err != nil {
            logsv.Debugln("TlsServer Listen error %v.", err.Error())
            break
        }
        logsv.Infof("TlsServer Accept Remote<%v> - Local<%v>.", conn.RemoteAddr().String(), conn.LocalAddr().String())
        t.bzn.HandleConn(conn)
    }
    logsv.Debugln("TlsServer Listen over..")
    return nil
}
