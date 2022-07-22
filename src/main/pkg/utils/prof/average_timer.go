package prof

import (
	"sync"
	"time"
)

// AverageTimer is a average with timer
type AverageTimer struct {
	lock     sync.RWMutex
	nowSum   int64
	nowCount int64
	nowMax   int64
	nowTimer int64
	lastMean float64
	lastMax  int64
	name     string
	period   int64
}

// NewAverageTimer return AverageTimer with params
func NewAverageTimer(name string, period int64) *AverageTimer {
	if period == 0 {
		period = 60
	}
	now := time.Now().Unix()
	tt := now - now%period

	return &AverageTimer{
		name:     name,
		period:   period,
		nowTimer: tt,

		nowSum:   0,
		nowCount: 0,
		nowMax:   0,
		lastMean: 0,
		lastMax:  0,
	}
}

// Set is to add value into Average
func (a *AverageTimer) Set(value int64) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.updatePeriod()
	if a.nowMax < value {
		a.nowMax = value
	}
	a.nowSum += value
	a.nowCount++
}

// Mean return the mean value of last data
func (a *AverageTimer) Mean() float64 {
	a.lock.RLock()
	defer a.lock.RUnlock()

	a.updatePeriod()
	return fixFloat(a.lastMean, 2)
}

// Max return the max value of last data
func (a *AverageTimer) Max() float64 {
	a.lock.RLock()
	defer a.lock.RUnlock()

	a.updatePeriod()
	return float64(a.lastMax)
}

// Name return the name of AverageTimer
func (a *AverageTimer) Name() string {
	return a.name
}

func (a *AverageTimer) updatePeriod() {
	now := time.Now().Unix()
	tt := now - now%a.period
	if tt != a.nowTimer {
		if a.nowCount != 0 {
			a.lastMean = float64(a.nowSum) / float64(a.nowCount)
			a.lastMax = a.nowMax
		} else {
			a.lastMean = 0
			a.lastMax = 0
		}

		a.nowTimer = tt
		a.nowSum = 0
		a.nowCount = 0
		a.nowMax = 0
	}
}
