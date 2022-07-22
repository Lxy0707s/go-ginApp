package prof

import (
	"sync"
	"time"
)

// CountBase is a simple counter
type CountBase struct {
	sync.RWMutex
	name string
	cnt  int64
	time int64
}

// NewCountBase return new counter with name
func NewCountBase(name string) *CountBase {
	uts := time.Now().Unix()
	return &CountBase{
		name: name,
		cnt:  0,
		time: uts}
}

// Cnt return counter value and counter time
func (cb *CountBase) Cnt() (int64, int64) {
	cb.RLock()
	defer cb.RUnlock()
	return cb.cnt, cb.time
}

// Name return counter's name
func (cb *CountBase) Name() string {
	cb.RLock()
	defer cb.RUnlock()
	return cb.name
}

// SetCnt sets counter value
func (cb *CountBase) SetCnt(cnt int64) {
	cb.Lock()
	cb.cnt = cnt
	cb.time = time.Now().Unix()
	cb.Unlock()
}

// Incr increase counter value by 1
func (cb *CountBase) Incr() {
	cb.IncrBy(1)
}

// IncrBy increase counter value by incr value
func (cb *CountBase) IncrBy(incr int64) {
	cb.Lock()
	cb.cnt += incr
	cb.time = time.Now().Unix()
	cb.Unlock()
}
