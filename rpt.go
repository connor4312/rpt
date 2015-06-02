package rpt

import (
	"time"
)

type RPT struct {
	periods    []uint
	resolution int64
	head       int64
	ptr        int
	size       int
	length     int
}

const (
	OVER_ALLOC = 4
)

// Creates a new rate tracker to hold `size` data points that are
// a given time long. For instance, to track requests per minute
// with one-second resolution, you can call New(60, time.Second)
func New(size int, resolution time.Duration) *RPT {
	return &RPT{
		periods:    make([]uint, size*OVER_ALLOC),
		resolution: resolution.Nanoseconds(),
		head:       0,
		ptr:        0,
		size:       size,
		length:     size * OVER_ALLOC,
	}
}

// Adds a single request in the current RPT time period. This is
// an alias to calling `AddRequestsTo(1, time.Now())`
func (r *RPT) AddRequest() {
	r.AddRequestsTo(1, time.Now())
}

// Adds the number of requests at the given time.
func (r *RPT) AddRequestsTo(requests uint, time time.Time) {
	now := time.UnixNano() / r.resolution
	target := int(now-r.head) + r.ptr

	// If we're pointing at something beyond the range of values,
	// we have to reshift the array.
	if target >= r.length {
		diff := target - r.length + 2
		if diff > r.size {
			diff = 0
		}

		r.shift(diff)
		target = diff
	}
	r.periods[target] += requests
	r.ptr = target
	r.head = now
}

// Shifts the data - wipes the array keeping the last n elements.
func (r *RPT) shift(n int) {
	for i := 0; i < r.length; i++ {
		if i < n {
			r.periods[i] = r.periods[r.length-n+i]
		} else {
			r.periods[i] = 0
		}
	}
}

// Returns the number of requests for the total time interval for
// which data has been recorded.
func (r *RPT) GetRPT() uint {
	var total uint
	min := r.ptr - r.size
	for i := r.ptr; i >= 0 && i > min; i-- {
		total += r.periods[i]
	}

	return total
}

// Returns an array of data in the range from start to end, based
// on the current time. E.g., GetRange(-60, 0) to get the last
// 60 recordings in chronological order.
func (r *RPT) GetRange(start, end int) []uint {
	start = r.filterBound(start)
	end = r.filterBound(end)

	// There was a smart bitwise hack to do this,
	// but I fail to recall it...
	var dir, size int
	if end > start {
		dir = 1
		size = end - start
	} else {
		dir = -1
		size = start - end
	}

	out := make([]uint, size+1)
	for i := 0; i <= size; i += 1 {
		out[i] = r.periods[r.ptr+start+i*dir]
	}

	return out
}

func (r *RPT) filterBound(bound int) int {
	if bound > 0 {
		return 0
	}
	if r.ptr+bound < 0 {
		return -r.ptr
	}

	return bound
}
