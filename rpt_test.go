package rpt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestZeroByDefault(t *testing.T) {
	r := New(60, time.Second)
	assert.Equal(t, uint(0), r.GetRPT())
}

func TestAddsContiguousDataPoints(t *testing.T) {
	r := New(60, time.Second)
	r.AddRequestsTo(1, time.Unix(1, 0))
	r.AddRequestsTo(2, time.Unix(2, 0))
	r.AddRequestsTo(3, time.Unix(3, 0))
	assert.Equal(t, uint(6), r.GetRPT())
}

func TestAddsNonContiguousDataPoints(t *testing.T) {
	r := New(60, time.Second)
	r.AddRequestsTo(1, time.Unix(1, 0))
	r.AddRequestsTo(2, time.Unix(20, 0))
	r.AddRequestsTo(3, time.Unix(30, 0))
	assert.Equal(t, uint(6), r.GetRPT())
}

func TestDiscardsOld(t *testing.T) {
	r := New(3, time.Second)
	r.AddRequestsTo(1, time.Unix(1, 0))
	r.AddRequestsTo(2, time.Unix(2, 0))
	r.AddRequestsTo(3, time.Unix(3, 0))
	assert.Equal(t, uint(6), r.GetRPT())
	// now it starts moving the pointer
	r.AddRequestsTo(4, time.Unix(4, 0))
	assert.Equal(t, uint(9), r.GetRPT())
	// continues moving
	r.AddRequestsTo(5, time.Unix(5, 0))
	assert.Equal(t, uint(12), r.GetRPT())
	// discards
	r.AddRequestsTo(8, time.Unix(7, 0))
	assert.Equal(t, uint(13), r.GetRPT())
	// adds after long interval
	r.AddRequestsTo(10, time.Unix(100, 0))
	assert.Equal(t, uint(10), r.GetRPT())
}

func TestRolloverBorders1(t *testing.T) {
	r := New(3, time.Second)
	// adds records at rollover borders
	k := 10
	r.AddRequestsTo(uint(8), time.Unix(int64(1), 0))
	r.AddRequestsTo(uint(9), time.Unix(int64(2), 0))
	for i := 3; i < 100; i++ {
		r.AddRequestsTo(uint(k), time.Unix(int64(i), 0))
		assert.Equal(t, uint(k*3-3), r.GetRPT())
		k++
	}
}

func TestRolloverBorders2(t *testing.T) {
	r := New(3, time.Second)
	// adds records at rollover borders
	k := 10
	r.AddRequestsTo(uint(9), time.Unix(int64(1), 0))
	for i := 3; i < 100; i += 2 {
		r.AddRequestsTo(uint(k), time.Unix(int64(i), 0))
		assert.Equal(t, uint(k*2-1), r.GetRPT())
		k++
	}
}

func TestRolloverBorders3(t *testing.T) {
	r := New(3, time.Second)
	// adds records at rollover borders
	k := 10
	for i := 0; i < 100; i += 3 {
		r.AddRequestsTo(uint(k), time.Unix(int64(i), 0))
		assert.Equal(t, uint(k), r.GetRPT())
		k++
	}
}

func TestRolloverBorders4(t *testing.T) {
	r := New(3, time.Second)
	// adds records at rollover borders
	k := 10
	for i := 0; i < 100; i += 4 {
		r.AddRequestsTo(uint(k), time.Unix(int64(i), 0))
		assert.Equal(t, uint(k), r.GetRPT())
		k++
	}
}

func TestGetsBounds(t *testing.T) {
	l := 5
	r := New(l, time.Second)
	for i := 0; i < l; i++ {
		r.AddRequestsTo(uint(i), time.Unix(int64(i), 0))
	}

	assert.Equal(t, []uint{4}, r.GetRange(0, 0))
	assert.Equal(t, []uint{3, 4}, r.GetRange(-1, 0))
	assert.Equal(t, []uint{0, 1, 2, 3, 4}, r.GetRange(-4, 0))
	assert.Equal(t, []uint{0, 1, 2, 3, 4}, r.GetRange(-500, 0))
	assert.Equal(t, []uint{4, 3, 2, 1, 0}, r.GetRange(0, -500))
	assert.Equal(t, []uint{1, 2, 3}, r.GetRange(-3, -1))
}

func BenchmarkAddRequest(b *testing.B) {
	r := New(60, time.Second)
	for n := 0; n < b.N; n++ {
		r.AddRequestsTo(uint(n), time.Unix(int64(n), 0))
	}
}

func BenchmarkGetRtp(b *testing.B) {
	l := 100
	r := New(l, time.Second)
	for i := 0; i < l; i++ {
		r.AddRequestsTo(uint(i), time.Unix(int64(i), 0))
	}

	for n := 0; n < b.N; n++ {
		r.GetRPT()
	}
}

func BenchmarkGetRange(b *testing.B) {
	l := 100
	r := New(l, time.Second)
	for i := 0; i < l; i++ {
		r.AddRequestsTo(uint(i), time.Unix(int64(i), 0))
	}

	for n := 0; n < b.N; n++ {
		r.GetRange(-90, -10)
	}
}
