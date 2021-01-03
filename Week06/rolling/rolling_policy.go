package rolling

import (
	"sync"
	"time"
)

type RollingPolicy struct {
	mu     sync.RWMutex
	size   int
	window *Window
	offset int

	bucketDuration time.Duration
	lastAppendTime time.Time
}

type RollingPolicyOpt struct {
	bucketDuration time.Duration
}

func NewRollingPolicy(w *Window, opt RollingPolicyOpt) *RollingPolicy {
	return &RollingPolicy{
		window:         w,
		size:           w.Size(),
		offset:         0,
		bucketDuration: opt.bucketDuration,
		lastAppendTime: time.Now(),
	}
}

func (r *RollingPolicy) timespan() int {
	v := int(time.Since(r.lastAppendTime) / r.bucketDuration)
	if v > -1 {
		return v
	}
	return r.size
}

func (r *RollingPolicy) add(f func(offset int, val float64), val float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	timespan := r.timespan()
	if timespan > 0 {
		r.lastAppendTime = r.lastAppendTime.Add(time.Duration(timespan * int(r.bucketDuration)))
		offset := r.offset
		s := offset + 1
		if timespan > r.size {
			timespan = r.size
		}
		e, e1 := s+timespan, 0
		if e > r.size {
			e1 = e - r.size
			e = r.size
		}
		for i := s; i < e; i++ {
			r.window.ResetBucket(i)
			offset = i
		}
		for i := 0; i < e1; i++ {
			r.window.ResetBucket(i)
			offset = i
		}
		r.offset = offset
	}
	f(r.offset, val)
}

func (r *RollingPolicy) Append(val float64) {
	r.add(r.window.Append, val)
}

func (r *RollingPolicy) Add(val float64) {
	r.add(r.window.Add, val)
}

func (r *RollingPolicy) Reduce(f func(Iterator) float64) (val float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	timespan := r.timespan()
	if count := r.size - timespan; count > 0 {
		offset := r.offset + timespan + 1
		if offset > r.size {
			offset = offset - r.size
		}
		val = f(r.window.Iterator(offset, count))
	}
	return val
}
