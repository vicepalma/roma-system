package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Limiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	r        rate.Limit
	burst    int
	ttl      time.Duration
}

func NewLimiter(r rate.Limit, burst int, ttl time.Duration) *Limiter {
	l := &Limiter{
		visitors: map[string]*visitor{},
		r:        r, burst: burst, ttl: ttl,
	}
	go l.cleanup()
	return l
}

func (l *Limiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.visitors[ip]
	if !ok {
		lim := rate.NewLimiter(l.r, l.burst)
		l.visitors[ip] = &visitor{limiter: lim, lastSeen: time.Now()}
		return lim
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (l *Limiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > l.ttl {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

func (l *Limiter) Gin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = "unknown"
		}
		lim := l.get(ip)
		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate_limited"})
			return
		}
		c.Next()
	}
}
