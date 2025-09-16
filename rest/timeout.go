package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// timeoutWare middleware that enforces a timeout.
// If the request takes longer than d, it returns 408 Request Timeout.
//
// Remember to use c.UserContext() in your handlers to get the context with timeout.
// c.UserContext() will return the original context if no timeout is set.
// c.Context() returns internal context of Fasthttp (Fiber uses it behind the scenes).
func timeoutWare(d time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a context with timeout based on the existing request context.
		ctx, cancel := context.WithTimeout(c.UserContext(), d)
		defer cancel()
		c.SetUserContext(ctx) // Attach the new context to Fiber's context.
		err := c.Next()       // Call the next handler in chain.
		// Check if the timeout expired.
		if err == nil && ctx.Err() == context.DeadlineExceeded {
			return errorJSON(c, http.StatusRequestTimeout, detailsResp{
				Message: "Request timeout",
				Details: "The server timed out waiting for the request.",
			})
		}
		return err
	}
}

// testTimeout godoc
//
//	@Summary		Test timeout middleware
//	@Description	Simulates 5 seconds of work
//	@Tags			debug
//	@Produce		json
//	@Failure		508	{object}	errorResp
//	@Success		200	{object}	resp{message=string}
//	@Router			/test-timeout [get]
func testTimeout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		time.Sleep(5 * time.Second)
		return respJSON(c, http.StatusOK, detailsResp{
			Message: "Finished work.",
		})
	}
}
