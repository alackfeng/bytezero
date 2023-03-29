package sysstat

import (
	"fmt"
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

// Execute -
func (s *SysStat) Execute() {
	fmt.Println(">>>>>bytezero sysstat Execute begin.")
	if err := s.netst.Stat(); err != nil {
		fmt.Println(">>>>>bytezero sysstat GetInterfaces err.", err.Error())
	}
	fmt.Println(">>>>>bytezero sysstat Execute end...")
}
