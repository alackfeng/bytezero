package sysstat

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// SysStat - 系统资源统计.
type SysStat struct {
	netst *NetStat
}

// NewSysStat -
func NewSysStat() *SysStat {
	return &SysStat{
		netst: NewNetStat(),
	}
}

// Init -
func (s *SysStat) Init() {
	if err := s.netst.Init(); err != nil {
		panic(err.Error())
	}
}

// GetAdapterList - 获取网卡适配器.
func (s *SysStat) GetAdapterList() (*syscall.IpAdapterInfo, error) {
	b := make([]byte, 1000)
	l := uint32(len(b))
	a := (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
	err := syscall.GetAdaptersInfo(a, &l)
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b = make([]byte, l)
		a = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
		err = syscall.GetAdaptersInfo(a, &l)
	}
	if err != nil {
		return nil, os.NewSyscallError("GetAdaptersInfo", err)
	}
	return a, nil
}

// Execute -
func (s *SysStat) Execute() {
	fmt.Println(">>>>>bytezero sysstat Execute begin.")
	if err := s.netst.Stat(); err != nil {
		fmt.Println(">>>>>bytezero sysstat GetInterfaces err.", err.Error())
	}
	fmt.Println(">>>>>bytezero sysstat Execute end...")
}
