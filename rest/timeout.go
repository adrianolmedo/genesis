package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// timeoutWare middleware that enforces a timeout.
// If the request takes longer than d, it returns 408 Request Timeout.
func timeoutWare(d time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a context with timeout based on the existing request context.
		ctx, cancel := context.WithTimeout(c.UserContext(), d)
		defer cancel()
		c.SetUserContext(ctx) // Attach the new context to Fiber's context.
		err := c.Next()       // Call the next handler in chain.
		// Check if the timeout expired.
		if ctx.Err() == context.DeadlineExceeded {
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
//	@Summary		Test timeout
//	@Description	Simulates a long-running operation to test context timeout
//	@Tags			debug
//	@Produce		json
//	@Failure		504	{object}	errorResp
//	@Success		200	{object}	resp{message=string}
//	@Router			/test-timeout [get]
func testTimeout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// We derive a context with a 2-second timeout from the request context.
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		done := make(chan string, 1)
		go func() {
			time.Sleep(5 * time.Second) // Simulate work that takes 5 seconds.
			done <- "Finished work."
		}()
		select {
		case <-ctx.Done(): // Timeout expired
			return errorJSON(c, http.StatusGatewayTimeout, detailsResp{
				Message: "The operation took too long",
				Details: "Please try again later.",
			})
		case result := <-done: // Response successful before timeout
			return respJSON(c, http.StatusOK, detailsResp{
				Message: result,
			})
		}
	}
}
