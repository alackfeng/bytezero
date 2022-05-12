package server

import (
	"fmt"
	"net"

	"github.com/alackfeng/bytezero/cores/utils"
)

var logsv = utils.Logger(utils.Fields{"animal": "server"})

// TcpServer -
type TcpServer struct {
    network string
    address string
    maxBufferLen int // recv max buffer len
}

// NewTcpServer -
func NewTcpServer(address string) *TcpServer {
    return &TcpServer{
        network: "tcp",
        address: address,
        maxBufferLen: 1024 * 1024 * 10,
    }
}

// Listen -
func (t *TcpServer) Listen() error {
    var err error = nil
    tcpAddr, err := net.ResolveTCPAddr(t.network, t.address)
    if err != nil {
        return err
    }
    tcpListener, err := net.ListenTCP(t.network, tcpAddr)
    if err != nil {
        return err
    }
    logsv.Debugln("TcpServer Listen begin.", tcpAddr.String())

    for {
        tcpConn, err := tcpListener.AcceptTCP()
        if err != nil {
            logsv.Debugln("TcpServer Listen error %v.", err.Error())
            break
        }
        logsv.Infof("TcpServer Accept Remote<%v> - Local<%v>.", tcpConn.RemoteAddr().String(), tcpConn.LocalAddr().String())
        go t.handleConn(tcpConn)

    }
    logsv.Debugln("TcpServer Listen over..")
    return err
}

// handleConn -
func (t *TcpServer) handleConn(conn *net.TCPConn) error {
    defer conn.Close()
    count := 0
    buffer := make([]byte, t.maxBufferLen)
    for {
        len, err := conn.Read(buffer)
        if err != nil {
            logsv.Errorf("TcpServer handleConn - read error.", err.Error())
            return err
        }
        if len == 0 {
            // logsv.Debugln("TcpServer handleConn - wait next.", count)
            count += 1
            continue
        }

        wlen, err := conn.Write(buffer[0:len])
        if err != nil {
            return err
        }
        fmt.Printf("No.%d: read %d, write %d.\n", count, len, wlen)
        count += 1
    }
}


