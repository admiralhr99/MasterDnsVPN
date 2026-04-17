package client

import (
	"sync"
	"time"
)

type sendRateLimiter struct {
	mu       sync.Mutex
	tokens   float64
	maxBurst float64
	rate     float64
	last     time.Time
}

func newSendRateLimiter(perSecond int) *sendRateLimiter {
	rate := float64(perSecond)
	burst := rate / 10
	if burst < 1 {
		burst = 1
	}
	return &sendRateLimiter{
		tokens:   burst,
		maxBurst: burst,
		rate:     rate,
		last:     time.Now(),
	}
}

func (l *sendRateLimiter) Take() {
	for {
		l.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(l.last).Seconds()
		l.last = now
		l.tokens += elapsed * l.rate
		if l.tokens > l.maxBurst {
			l.tokens = l.maxBurst
		}
		if l.tokens >= 1.0 {
			l.tokens -= 1.0
			l.mu.Unlock()
			return
		}
		deficit := 1.0 - l.tokens
		wait := time.Duration(deficit / l.rate * float64(time.Second))
		l.mu.Unlock()
		time.Sleep(wait)
	}
}
