package prof

import (
	"sync"
	"time"
)

// CountPeriod is a simple counter with period
type CountPeriod struct {
	sync.RWMutex
	name string

	period     int64
	cntMap     map[int64]int64
	keepPeriod int64
}

// NewCountPeriod return new counter with name
func NewCountPeriod(name string, period int64, keepPeriod int64) *CountPeriod {
	return &CountPeriod{
		name:       name,
		period:     period,
		cntMap:     make(map[int64]int64),
		keepPeriod: keepPeriod,
	}
}

// Cnt return counter value and counter time
func (cp *CountPeriod) Cnt() (int64, int64) {
	cp.RLock()
	defer cp.RUnlock()
	uts := time.Now().Unix()
	uts = uts - uts%cp.period - cp.period
	return cp.cntMap[uts], uts
}

// Name return counter's name
func (cp *CountPeriod) Name() string {
	cp.RLock()
	defer cp.RUnlock()
	return cp.name
}

// SetCnt sets counter value
func (cp *CountPeriod) SetCnt(cnt int64) {
	cp.Lock()
	uts := time.Now().Unix()
	uts = uts - uts%cp.period
	cp.cntMap[uts] = cnt
	cp.Unlock()
}

// Incr increase counter value by 1
func (cp *CountPeriod) Incr() {
	cp.IncrBy(1)
}

// IncrBy increase counter value by incr value
func (cp *CountPeriod) IncrBy(incr int64) {
	cp.Lock()
	uts := time.Now().Unix()
	uts = uts - uts%cp.period
	cp.cntMap[uts] = cp.cntMap[uts] + incr
	cp.clear(uts)
	cp.Unlock()
}

// GetAll return all data in map, order by time desc
func (cp *CountPeriod) GetAll() []int64 {
	cp.RLock()
	defer cp.RUnlock()
	uts := time.Now().Unix()
	uts = uts - uts%cp.period - cp.period
	result := make([]int64, 0)
	for tt := uts; tt >= uts-cp.keepPeriod*cp.period; tt -= cp.period {
		result = append(result, cp.cntMap[tt])
	}
	return result
}

func (cp *CountPeriod) clear(now int64) {
	overtime := now - cp.keepPeriod*cp.period
	for tt := range cp.cntMap {
		if tt < overtime {
			delete(cp.cntMap, tt)
		}
	}
}
