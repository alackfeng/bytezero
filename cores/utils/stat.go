package utils

import (
	"fmt"
	"time"
)

// StatBandwidth - 统计带宽bps.
type StatBandwidth struct {
    BeginTime time.Time
    EndTime time.Time
    Count int64
    Bytes int64
}

// Begin -
func (s *StatBandwidth) Begin() {
    s.BeginTime = time.Now()
}

// End -
func (s *StatBandwidth) End() {
    s.EndTime = time.Now()
}

// Inc -
func (s *StatBandwidth) Inc(b int64) {
    s.Count += 1
    s.Bytes += b
}

// String -
func (s *StatBandwidth) String() string {
    return fmt.Sprintf("(time: %v => %v) - %v bytes(%v count)", s.BeginTime.Format("2006-01-02 15:04:05.999999999"), s.EndTime.Format("2006-01-02 15:04:05.999999999"), s.Bytes, s.Count)
}

// Info -
func (s *StatBandwidth) Info() string {
    return fmt.Sprintf("(time: %v => %v) - %v bytes(%v count)", s.BeginTime.Format("2006-01-02 15:04:05.999999999"), s.EndTime.Format("2006-01-02 15:04:05.999999999"), s.Bytes, s.Count)
}

func (s *StatBandwidth) InfoAll() string {
    return fmt.Sprintf("(time: %s => %s<dura:%d ms>) - %d bytes(%d count)", s.BeginTime.Format("2006-01-02 15:04:05.999999999"), s.EndTime.Format("2006-01-02 15:04:05.999999999"), s.EndTime.Sub(s.BeginTime), s.Bytes, s.Count)
}
