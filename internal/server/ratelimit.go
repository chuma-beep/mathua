package server

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	burst    int
	interval time.Duration
}

type visitor struct {
	tokens    int
	lastCheck time.Time
}

func newRateLimiter(rate, burst int, interval time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		interval: interval,
	}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-30 * time.Minute)
		for ip, v := range rl.visitors {
			if v.lastCheck.Before(cutoff) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[ip]
	now := time.Now()

	if !ok {
		rl.visitors[ip] = &visitor{tokens: rl.burst - 1, lastCheck: now}
		return true
	}

	elapsed := now.Sub(v.lastCheck)
	v.lastCheck = now
	v.tokens += int(elapsed / rl.interval)
	if v.tokens > rl.burst {
		v.tokens = rl.burst
	}
	if v.tokens <= 0 {
		return false
	}
	v.tokens--
	return true
}

func (rl *rateLimiter) middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}
		if !rl.allow(ip) {
			http.Error(w, `{"error":"rate limit exceeded"}`, 429)
			return
		}
		next(w, r)
	}
}
