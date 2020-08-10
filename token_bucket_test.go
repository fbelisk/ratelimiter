package rateLimiter

import (
	"fmt"
	"testing"
	"time"
)

func TestBucket_Put(t *testing.T) {
	type fields struct {
		closeCh  chan int
		Capacity int64
		Tokens   int64
		Interval time.Duration
		Inc      int64
	}
	type args struct {
		Inc int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{
				Capacity: 1000,
				Tokens:   0,
				Interval: 0,
				Inc:      0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bucket{
				closeCh:  tt.fields.closeCh,
				Capacity: tt.fields.Capacity,
				tokens:   tt.fields.Tokens,
				Interval: tt.fields.Interval,
				Inc:      tt.fields.Inc,
			}
			for i:= 0; i < 500; i++ {
				go b.Put(1)
			}
			//for i:= 0; i < 499; i++ {
			//	go b.Take(0)
			//}
			fmt.Printf("%+v", b)
		})
	}
}

func TestBucket_TakeWait(t *testing.T) {
	type fields struct {
		closeCh  chan int
		Capacity int64
		tokens   int64
		Interval time.Duration
		Inc      int64
	}
	type args struct {
		count   int64
		maxWait time.Duration
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantTokens   int64
		wantWaitTime time.Duration
	}{
		{
			name: "base test",
			args:args{5, 6*time.Second},
			wantTokens:5,
			wantWaitTime:5*time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(10, 1, time.Second)
			gotTokens, gotWaitTime := b.TakeWait(tt.args.count, tt.args.maxWait)
			if gotTokens != tt.wantTokens {
				t.Errorf("TakeWait() gotTokens = %v, want %v", gotTokens, tt.wantTokens)
			}
			if gotWaitTime != tt.wantWaitTime {
				t.Errorf("TakeWait() gotWaitTime = %v, want %v", gotWaitTime, tt.wantWaitTime)
			}
		})
	}
}

func TestBucket_New(t *testing.T) {
	t.Run("base_test", func(t *testing.T) {
		New(10, 1, time.Second)
	})
}