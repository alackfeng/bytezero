package server

import (
	"crypto/rand"
	"crypto/tls"
	"net"
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

// Serve -
func (t *TlsServer) Serve(l net.Listener) (err error) {
    for {
        conn, err := l.(*net.TCPListener).AcceptTCP()
        if err != nil {
            logsv.Debugln("TlsServer Listen error %v.", err.Error())
            break
        }
        logsv.Infof("TlsServer Accept Remote<%v> - Local<%v>.", conn.RemoteAddr().String(), conn.LocalAddr().String())
        t.bzn.HandleConn(conn)
    }
    return nil
}

// Start -
func (t *TlsServer) Start() error {
    return bz.ListenTLS(t.address, t.certFile, t.keyFile, t)
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
