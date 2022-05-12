package client

import "net"

// TcpClient -
type TcpClient struct {
    *net.TCPConn
    network string
    address string
}

// NewTcpClient -
func NewTcpClient(address string) *TcpClient {
    return &TcpClient{
        network: "tcp",
        address: address,
    }
}

// Dial -
func (t *TcpClient) Dial() error {
    tcpAddr, err := net.ResolveTCPAddr(t.network, t.address)
    if err != nil {
        return err
    }
    tcpConn, err := net.DialTCP(t.network, nil, tcpAddr)
    if err != nil {
        return err
    }
    t.TCPConn = tcpConn
    return nil
}


