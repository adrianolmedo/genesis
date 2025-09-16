package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// testTimeout godoc
//
//	@Summary		Test timeout of 2 seconds
//	@Description	Simulates 5 seconds of work
//	@Tags			debug
//	@Produce		json
//	@Failure		408	{object}	errorResp
//	@Success		200	{object}	resp{message=string}
//	@Router			/test-timeout [get]
func testTimeout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		select {
		case <-time.After(5 * time.Second): // Simulate work
			return respJSON(c, http.StatusOK, detailsResp{
				Message: "Finished work.",
			})
		case <-ctx.Done(): // Timeout or cancellation
			return errorJSON(c, http.StatusRequestTimeout, detailsResp{
				Message: "Request timeout",
				Details: "The server timed out waiting for the request.",
			})
		}
	}
}
