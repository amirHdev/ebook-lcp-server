package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	mu        sync.Mutex
	limit     int
	window    time.Duration
	requests  map[string]int
	resetTime time.Time
}

func New(limit int, window time.Duration) *Limiter {
	return &Limiter{
		limit:     limit,
		window:    window,
		requests:  map[string]int{},
		resetTime: time.Now().Add(window),
	}
}

func (l *Limiter) Allow(key string) bool {
	if l == nil || l.limit <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if now.After(l.resetTime) {
		l.requests = map[string]int{}
		l.resetTime = now.Add(l.window)
	}
	l.requests[key]++
	return l.requests[key] <= l.limit
}
