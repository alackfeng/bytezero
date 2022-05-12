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
    // bps.
    nextTime time.Time
    nextBytes int64
    lastBps int64
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
    s.nextTime = time.Now()
    s.nextBytes += b
}

// Bps -
func (s *StatBandwidth) Bps(ms int64) int64 {
    dura := time.Now().Sub(s.nextTime)
    if dura.Milliseconds() > ms {
        s.lastBps = s.nextBytes * ms / dura.Milliseconds();
        s.nextTime = time.Now()
        s.nextBytes = 0
    }
    return s.lastBps
}

// Bps1s -
func (s *StatBandwidth) Bps1s() int64 {
    return s.Bps(1000)
}

// Bps3s -
func (s *StatBandwidth) Bps3s() int64 {
    return s.Bps(3000)
}

// Bps5s -
func (s *StatBandwidth) Bps5s() int64 {
    return s.Bps(5000)
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
