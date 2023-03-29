package sysstat

import (
	"net"
	"syscall"
)

// GetNetworkMaxSpeed - 获取网卡最大速率.
func GetNetworkMaxSpeed(networkInterface net.Interface) (uint32, error) {
	pIfRow := &syscall.MibIfRow{Index: uint32(networkInterface.Index)}
	if err := syscall.GetIfEntry(pIfRow); err != nil {
		return 0, err
	}
	return pIfRow.Speed, nil
}
