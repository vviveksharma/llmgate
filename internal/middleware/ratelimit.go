package middleware

import (
	"net/http"

	"github.com/vviveksharma/llmgate/internal/ratelimit"
)

type RateLimit struct {
	bucket *ratelimit.TokenBucket
}

func NewRateLimit(rate float64, burst int) RateLimit {
	return RateLimit{bucket: ratelimit.NewTokenBucket(rate, burst)}
}

func (rl RateLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.bucket.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
