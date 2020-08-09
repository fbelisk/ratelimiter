package rateLimiter

import (
	"sync/atomic"
	"time"
)

type Bucket struct {
	closeCh  chan int
	Capacity int64
	Tokens   int64
	Interval time.Duration
	Inc      int64
}

func (b *Bucket) Start() {
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
		t := atomic.LoadInt64(&b.Tokens)
		if t+count > b.Capacity {
			if atomic.CompareAndSwapInt64(&b.Tokens, t, b.Capacity) {
				return b.Capacity - t
			}
			continue
		}
		if atomic.CompareAndSwapInt64(&b.Tokens, t, t+count) {
			return count
		}
		continue
	}
}

func (b *Bucket) Take(count int64) int64 {
	if count <= 0 {
		return 0
	}
	for {
		t := atomic.LoadInt64(&b.Tokens)
		if t < count {
			if atomic.CompareAndSwapInt64(&b.Tokens, t, 0) {
				return t
			}
			continue
		}

		if atomic.CompareAndSwapInt64(&b.Tokens, t, t-count) {
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
		t := atomic.LoadInt64(&b.Tokens)
		if t < count {
			waitTime = time.Duration((count - t + b.Inc) * int64(b.Interval) / b.Inc)
			if waitTime > maxWait {
				waitTime = maxWait
				count = int64(maxWait /b.Interval) * b.Inc + t - b.Inc
			}
		}

		if atomic.CompareAndSwapInt64(&b.Tokens, t, t-count) {
			time.Sleep(waitTime)
			return count, waitTime
		}
		continue
	}
}
