package network

import (
	"golang.org/x/time/rate"
	"net/http"
)

const requestsPerSecond = 100

var limiter = rate.NewLimiter(requestsPerSecond, 1)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}