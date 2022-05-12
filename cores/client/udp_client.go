package client

import "net"

// UdpClient -
type UdpClient struct {
    *net.UDPConn
    network string
    address string
}

// NewUdpClient -
func NewUdpClient(address string) *UdpClient {
    return &UdpClient{
        network: "udp",
        address: address,
    }
}

// Dial -
func (u *UdpClient) Dial() error {
    udpAddr, err := net.ResolveUDPAddr(u.network, u.address)
    if err != nil {
        return err
    }
    udpConn, err := net.DialUDP(u.network, nil, udpAddr)
    if err != nil {
        return err
    }
    u.UDPConn = udpConn
    return nil
}
