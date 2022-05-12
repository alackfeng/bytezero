package server

import "net"

// UdpServer -
type UdpServer struct {
    network string
    address string
    maxBufferLen int
}

// NewUdpServer -
func NewUdpServer(address string) *UdpServer {
    return &UdpServer{
        network: "udp",
        address: address,
        maxBufferLen: 1024 * 1024 * 10,
    }
}

func (u *UdpServer) Listen() error {
    var err error = nil
    udpAddr, err := net.ResolveUDPAddr(u.network, u.address)
    if err != nil {
        return err
    }
    udpListener, err := net.ListenUDP(u.network, udpAddr)
    if err != nil {
        return err
    }
    logsv.Debugln("UdpServer Listen begin.", udpAddr.String())

    buffer := make([]byte, u.maxBufferLen)
    for {
        n, addr, err := udpListener.ReadFromUDP(buffer[:])
        if err != nil {
            logsv.Debugln("UdpServer Listen error %v.", err.Error())
            break
        }
        logsv.Infof("UdpServer Accept Remote<%v> - msg %d.", addr.String(), n)
        r, err := udpListener.WriteToUDP(buffer[0:n], addr)
        if err != nil {
            logsv.Infof("UdpServer WriteToUDP Remote<%v> - errror.%v", addr.String(), err.Error())
            continue
        }
        if r != n {
            logsv.Infof("UdpServer WriteToUDP Remote<%v> - wrong len %d, real %d.", addr.String(), r, n)
        }
    }
    logsv.Debugln("UdpServer Listen over..")
    return err
}


