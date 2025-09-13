package rest

import (
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

// rateLimit struct to hold rate limiting data.
type rateLimit struct {
	limiters sync.Map
	rps      rate.Limit
	burst    int
	ttl      time.Duration
}

// newRateLimit creates a new rateLimit instance.
func newRateLimit(rps rate.Limit, burst int, ttl time.Duration) *rateLimit {
	rl := &rateLimit{
		rps:   rps,
		burst: burst,
		ttl:   ttl,
	}
	// start cleanup goroutine to remove old limiters
	go rl.cleanupWorker()
	return rl
}

// limiterEntry holds the limiter and the last seen time.
type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// limiterPerIP retrieves or creates a rate limiter for a given IP.
func (rl *rateLimit) limiterPerIP(ip string) *rate.Limiter {
	val, ok := rl.limiters.Load(ip)
	if ok {
		entry := val.(*limiterEntry)
		entry.lastSeen = time.Now()
		return entry.limiter
	}
	limiter := rate.NewLimiter(rl.rps, rl.burst)
	rl.limiters.Store(ip, &limiterEntry{
		limiter:  limiter,
		lastSeen: time.Now(),
	})
	return limiter
}

// middlewarePerIP is a Fiber middleware that limits requests per IP.
func (rl *rateLimit) middlewarePerIP(c *fiber.Ctx) error {
	limiter := rl.limiterPerIP(c.IP())
	if !limiter.Allow() {
		return errorJSON(c, http.StatusTooManyRequests, detailsResp{
			Message: "Too many requests",
			Details: "You have exceeded the request limit. Please try again later.",
		})
	}
	return c.Next()
}

// cleanupWorker periodically removes old limiters.
func (rl *rateLimit) cleanupWorker() {
	ticker := time.NewTicker(rl.ttl)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		rl.limiters.Range(func(key, value any) bool {
			entry := value.(*limiterEntry)
			if now.Sub(entry.lastSeen) > rl.ttl {
				rl.limiters.Delete(key)
			}
			return true
		})
	}
}

// testRatelimit godoc
//
//	@Summary		Test rate limit
//	@Description	Simulates a long-running operation to test rate limit
//	@Tags			debug
//	@Produce		json
//	@Success		200	{object}	resp{message=string}
//	@Router			/test-ratelimit [get]
func testRatelimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// simulate processing
		time.Sleep(100 * time.Millisecond)
		return respJSON(c, http.StatusOK, detailsResp{
			Message: "You are within the rate limit.",
		})
	}
}
