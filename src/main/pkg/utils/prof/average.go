package prof

import (
	"sync"
)

// Average is a simple counter with average
type Average struct {
	lock sync.RWMutex
	list []int64
	size int
	name string
}

// NewAverage return new Average with params
func NewAverage(name string, size int) *Average {
	return &Average{
		name: name,
		size: size,
		list: make([]int64, 0, size*3),
	}
}

// Set is to add value into Average
func (a *Average) Set(value int64) {
	a.lock.Lock()
	a.list = append([]int64{value}, a.list...)
	if len(a.list) > a.size {
		a.list = a.list[:a.size]
	}
	a.lock.Unlock()
}

// Mean return the mean value of now data
func (a *Average) Mean() float64 {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if len(a.list) == 0 {
		return 0
	}
	var sum float64
	for _, v := range a.list {
		sum += float64(v)
	}
	return fixFloat(sum/float64(len(a.list)), 2)
}

// Name return the name of Average
func (a *Average) Name() string {
	return a.name
}
