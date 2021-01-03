package rolling

import (
	"fmt"
	"time"
)

type RollingCounterOpts struct {
	Size           int
	BucketDuration time.Duration
}

type RollingCounter struct {
	policy *RollingPolicy
}

func NewRollingCounter(opt RollingCounterOpts) *RollingCounter {
	window := NewWindow(WindowOpt{Size: opt.Size})
	policy := NewRollingPolicy(window, RollingPolicyOpt{bucketDuration: opt.BucketDuration})
	return &RollingCounter{policy: policy}
}

func (r *RollingCounter) Add(val int64) {
	if val < 0 {
		panic(fmt.Errorf("rolling: cannot decrease in value. val: %d", val))
	}
	r.policy.Add(float64(val))
}

func (r *RollingCounter) Reduce(f func(Iterator) float64) float64 {
	return r.policy.Reduce(f)
}

func (r *RollingCounter) Avg() float64 {
	return r.policy.Reduce(Avg)
}

func (r *RollingCounter) Min() float64 {
	return r.policy.Reduce(Min)
}

func (r *RollingCounter) Max() float64 {
	return r.policy.Reduce(Max)
}

func (r *RollingCounter) Sum() float64 {
	return r.policy.Reduce(Sum)
}

func (r *RollingCounter) Value() int64 {
	return int64(r.Sum())
}

func (r *RollingCounter) TimeSpan() int {
	return r.policy.timespan()
}
