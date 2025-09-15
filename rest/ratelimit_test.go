package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestRateLimit(t *testing.T) {
	app := fiber.New()
	rl := newRateLimit(2, 2, time.Second) // allow 2 req/sec, burst=2, ttl=1s
	app.Use(rl.middlewarePerIP)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected: 200, got: %d", res.StatusCode)
	}
	res, _ = app.Test(req)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected: 200, got: %d", res.StatusCode)
	}
	res, _ = app.Test(req)
	if res.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected: 429, got: %d", res.StatusCode)
	}
}

func TestRateLimitPerIP(t *testing.T) {
	app := fiber.New()
	rl := newRateLimit(1, 1, time.Second) // allow 1 req/sec, burst=1, ttl=1s
	app.Use(rl.middlewarePerIP)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	// IP1 request -> allowed
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = "10.0.0.1:1111"
	resp, _ := app.Test(req1)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected: 200, got: %d", resp.StatusCode)
	}
	// Same IP immediately again -> blocked
	resp, _ = app.Test(req1)
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected: 429, got: %d", resp.StatusCode)
	}
	// Wait for 1 second to reset the rate limiter
	time.Sleep(1 * time.Second)
	// Different IP should still pass
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.RemoteAddr = "10.0.0.2:2222"
	resp, _ = app.Test(req2)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected: 200, got: %d", resp.StatusCode)
	}
}

func TestRateLimitCleanup(t *testing.T) {
	rl := newRateLimit(1, 1, 50*time.Millisecond) // Very short TTL
	lim := rl.perIP("1.2.3.4")
	if lim == nil {
		t.Fatal("expected limiter")
	}
	time.Sleep(120 * time.Millisecond) // Wait longer than TTL so cleanup can remove it
	// Trigger cleanup manually
	rl.limiters.Range(func(key, _ any) bool {
		t.Fatalf("expected no entries after cleanup, found key %v", key)
		return true
	})
}
