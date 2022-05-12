package utils

import (
	"fmt"
	"time"
)

// Duration -
type Duration struct {
	curr int64 // 统计时间点.
}

func NewDuration() *Duration {
	return &Duration{
		curr: time.Now().UnixNano(),
	}
}

// DuraMs - Ms间隔.
func (d *Duration) DuraMs() int64 {
	return (time.Now().UnixNano() - d.curr) / 1e6
}

// DuraS - 秒间隔.
func (d *Duration) DuraS() int64 {
	return (time.Now().UnixNano() - d.curr) / 1e9
}

// DuraNano - 纳秒间隔.
func (d *Duration) DuraNano() int64 {
	return (time.Now().UnixNano() - d.curr)
}

// DuraNano - 纳秒间隔.
func (d *Duration) String() string {
	return fmt.Sprintf("%d Nano", (time.Now().UnixNano() - d.curr))
}
