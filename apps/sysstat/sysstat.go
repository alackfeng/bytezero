package sysstat

import (
	"context"
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
func (s *SysStat) Execute() context.CancelFunc {
	fmt.Println(">>>>>bytezero sysstat Execute begin.")
	ctx, cancel := context.WithCancel(context.TODO())
	if err := s.netst.Stat(ctx); err != nil {
		fmt.Println(">>>>>bytezero sysstat GetInterfaces err.", err.Error())
	}
	fmt.Println(">>>>>bytezero sysstat Execute end...")
	return cancel
}
