package sysstat

import (
	"fmt"
	"net"
	"syscall"
	"time"

	gnet "github.com/shirou/gopsutil/net"
)

// NetStat - 网卡采集.
type NetStat struct {
	netif        []NetInfo
	countersStat map[int64][]gnet.IOCountersStat
}

// NetInfo - 采集信息.
type NetInfo struct {
	net.Interface
	Addr  []net.Addr
	Speed uint32 // 网卡带宽最大速率.
}

func (s NetInfo) String() string {
	return fmt.Sprintf("name:%s - mac:%s - addrs:%v - speed:%v", s.Name, s.HardwareAddr.String(), s.Addr, s.Speed)
}

// NewNetStat -
func NewNetStat() *NetStat {
	return &NetStat{
		countersStat: make(map[int64][]gnet.IOCountersStat),
	}
}

func (s *NetStat) Init() error {
	return s.GetInterfaces()
}

// GetInterfaces - 获取网卡信息.
func (s *NetStat) GetInterfaces() error {
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for i, networkInterface := range networkInterfaces {
		networkInterfacei := networkInterface
		addr, _ := networkInterface.Addrs()
		pIfRow := &syscall.MibIfRow{Index: uint32(networkInterfacei.Index)}
		if err := syscall.GetIfEntry(pIfRow); err != nil {
			return err
		}
		fmt.Println(">>>>>bytezero sysstat interface ", i, networkInterfacei, addr[1], pIfRow.Speed/1024/1024)
		s.netif = append(s.netif, NetInfo{
			Interface: networkInterface,
			Addr:      addr,
			Speed:     pIfRow.Speed,
		})
	}
	return nil
}

// Stat -
func (s *NetStat) Stat() error {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixMilli()
			ioCounters, err := gnet.IOCounters(true)
			if err != nil {
				return err
			}
			s.countersStat[now] = ioCounters
			for i, ioCounter := range ioCounters {
				fmt.Println(">>>>>bytezero sysstat ioCounter ", i, ioCounter.Name, ioCounter.BytesSent, ioCounter.BytesRecv)
			}
		}
	}
}
