package server

import (
	"fmt"
	"net"

	bz "github.com/alackfeng/bytezero/bytezero"
	"github.com/alackfeng/bytezero/cores/utils"
)

var logsv = utils.Logger(utils.Fields{"animal": "server"})

// TcpServer -
type TcpServer struct {
    bzn bz.BZNet
    network string
    address string
    maxBufferLen int // recv max buffer len
    readBufferLen int
    writeBufferLen int
}

// NewTcpServer -
func NewTcpServer(bzn bz.BZNet, address string, maxBufferLen int, rwBufferLen int) *TcpServer {
    return &TcpServer{
        bzn: bzn,
        network: "tcp",
        address: address,
        maxBufferLen: maxBufferLen,
        readBufferLen: rwBufferLen,
        writeBufferLen: rwBufferLen,
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
        t.bzn.HandleConn(tcpConn)
    }
    logsv.Debugln("TcpServer Listen over..")
    return err
}

// handleEcho -
func (t *TcpServer) handleEcho(conn *net.TCPConn) error {
    // if t.readBufferLen > 1024 {
    //     conn.SetReadBuffer(t.readBufferLen)
    // }
    // if t.writeBufferLen > 1024 {
    //     conn.SetWriteBuffer(t.writeBufferLen)
    // }

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
        if count % 1000 == 0 {
            fmt.Printf("No.%d: read %d, write %d.\n", count, len, wlen)
        }
        count += 1
    }
}


