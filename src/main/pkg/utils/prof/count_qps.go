package prof

import (
	"sync"
	"time"
)

// CountQPS is a counter to save value and qps
type CountQPS struct {
	sync.RWMutex
	name     string
	cnt      int64
	qps      float64
	time     int64
	lastTime int64
	lastCnt  int64
}

// NewCountQPS return qps counter with name
func NewCountQPS(name string) *CountQPS {
	now := time.Now().Unix()
	return &CountQPS{
		name:     name,
		cnt:      0,
		qps:      0,
		time:     now,
		lastTime: now,
		lastCnt:  0}
}

// Cnt return counter value and time
func (cq *CountQPS) Cnt() (int64, int64) {
	cq.RLock()
	defer cq.RUnlock()
	return cq.cnt, cq.time
}

// QPS return counter qps value
func (cq *CountQPS) QPS() float64 {
	cq.Lock()
	defer cq.Unlock()

	now := time.Now().Unix()
	duration := now - cq.lastTime
	if duration >= 10 { // every 10 seconds to calculate qps value
		cq.qps = fixFloat(float64(cq.cnt-cq.lastCnt)/float64(duration), 2)
		cq.lastCnt = cq.cnt
		cq.lastTime = now
		return cq.qps
	}
	return cq.qps
}

// Name return counter's name
func (cq *CountQPS) Name() string {
	cq.RLock()
	defer cq.RUnlock()
	return cq.name
}

// Incr increase couter by 1
func (cq *CountQPS) Incr() {
	cq.IncrBy(1)
}

// IncrBy increase counter by incr
func (cq *CountQPS) IncrBy(incr int64) {
	cq.Lock()
	cq.cnt += incr
	cq.time = time.Now().Unix()
	cq.Unlock()
}
