package rateLimiter

import (
	"sync/atomic"
	"time"
)

type Bucket struct {
	closeCh  chan int
	Capacity int64
	tokens   int64
	Interval time.Duration
	Inc      int64
}

func New(capacity, inc int64, interval time.Duration) *Bucket{
	b := &Bucket{
		closeCh:  make(chan int),
		Capacity: capacity,
		Interval: interval,
		Inc:      inc,
	}
	go b.Run()

	return b
}

func (b *Bucket) Run() {
	ticker := time.NewTicker(b.Interval)
	defer ticker.Stop()

	for _ = range ticker.C {
		select {
		case <-b.closeCh:
			return
		default:
			b.Put(b.Inc)
		}
	}
}

func (b *Bucket) Put(count int64) int64 {
	if count <= 0 {
		return 0
	}
	for {
		t := atomic.LoadInt64(&b.tokens)
		if t == b.Capacity {
			return 0
		}
		if t+count > b.Capacity {
			if atomic.CompareAndSwapInt64(&b.tokens, t, b.Capacity) {
				return b.Capacity - t
			}
			continue
		}
		if atomic.CompareAndSwapInt64(&b.tokens, t, t+count) {
			return count
		}
	}
}

func (b *Bucket) Take(count int64) int64 {
	if count <= 0 {
		return 0
	}
	for {
		t := atomic.LoadInt64(&b.tokens)
		if t < count {
			if atomic.CompareAndSwapInt64(&b.tokens, t, 0) {
				return t
			}
			continue
		}

		if atomic.CompareAndSwapInt64(&b.tokens, t, t-count) {
			return count
		}
		continue
	}
}

func (b *Bucket) TakeWait(count int64, maxWait time.Duration) (tokens int64, waitTime time.Duration) {
	if count <= 0 {
		return 0, 0
	}
	for {
		t := atomic.LoadInt64(&b.tokens)
		if t < count {
			waitTime = time.Duration((count - t) * int64(b.Interval) / b.Inc + 1)
			if waitTime > maxWait {
				waitTime = maxWait
				count = int64(maxWait / b.Interval) * b.Inc + t
			}
		}
		if atomic.CompareAndSwapInt64(&b.tokens, t, t-count) {
			time.Sleep(waitTime)
			return count, waitTime
		}
	}
}
