package sysstat

import (
	"context"
	"fmt"
	"net"
	"time"

	gnet "github.com/shirou/gopsutil/net"
)

// NetStat - 网卡采集.
type NetStat struct {
	netif        []NetInfo
	countersStat []*CountersStat
	maxStat      int32
}

type CountersStat struct {
	Stat  []gnet.IOCountersStat `form:"Stat" json:"Stat" xml:"Stat" bson:"Stat" binding:"required"`
	NowMs int64                 `form:"nowMs" json:"nowMs" xml:"nowMs" bson:"nowMs" binding:"required"`
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

// StatNetResult -
type StatNetResult struct {
	Name  string          `form:"Name" json:"Name" xml:"Name" bson:"Name" binding:"required"`
	Mac   string          `form:"Mac" json:"Mac" xml:"Mac" bson:"Mac" binding:"required"`
	IP    []string        `form:"IP" json:"IP" xml:"IP" bson:"IP" binding:"required"`
	Speed uint32          `form:"Speed" json:"Speed" xml:"Speed" bson:"Speed" binding:"required"`
	Stats []*CountersStat `form:"Stats" json:"Stats" xml:"Stats" bson:"Stats" binding:"required"`
}

type StatNetResults struct {
	Info []StatNetResult `form:"Info" json:"Info" xml:"Info" bson:"Info" binding:"required"`
}

func (s *SysStat) GetNet() (nets StatNetResults, err error) {

	for _, net := range s.netst.netif {
		netResult := StatNetResult{
			Name:  net.Name,
			Mac:   net.HardwareAddr.String(),
			IP:    []string{},
			Speed: net.Speed,
		}
		addrs, err := net.Addrs()
		if err != nil {
			return nets, err
		}
		for _, addr := range addrs {
			netResult.IP = append(netResult.IP, addr.String())
		}
		netResult.Stats = s.netst.countersStat
		nets.Info = append(nets.Info, netResult)
	}
	return nets, err
}

// NewNetStat -
func NewNetStat() *NetStat {
	return &NetStat{
		maxStat: 20,
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
	for _, networkInterface := range networkInterfaces {
		var speed uint32
		addr, _ := networkInterface.Addrs()
		speed, err := GetNetworkMaxSpeed(networkInterface)
		if err != nil {
			return err
		}
		s.netif = append(s.netif, NetInfo{
			Interface: networkInterface,
			Addr:      addr,
			Speed:     speed,
		})
	}
	return nil
}

// Stat -
func (s *NetStat) Stat(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ctx.Done():
			fmt.Println(">>>>>bytezero sysstat ioCounter done.")
			break
		case <-ticker.C:
			now := time.Now().UnixMilli()
			ioCounters, err := gnet.IOCounters(true)
			if err != nil {
				return err
			}
			s.countersStat = append(s.countersStat, &CountersStat{Stat: ioCounters, NowMs: now})
			if len(s.countersStat) > int(s.maxStat) {
				s.countersStat = s.countersStat[s.maxStat/3:]
			}
			// for i, ioCounter := range ioCounters {
			// 	fmt.Println(">>>>>bytezero sysstat ioCounter ", i, ioCounter.Name, ioCounter.BytesSent, ioCounter.BytesRecv)
			// }
		}
	}
}
