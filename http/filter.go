package http

import (
	"fmt"
	"strings"

	domain "github.com/adrianolmedo/genesis"

	"github.com/gofiber/fiber/v2"
)

func getFilter(c *fiber.Ctx) (domain.Filter, error) {
	limit := c.QueryInt("limit")
	if limit < 0 {
		return domain.Filter{}, fmt.Errorf("positive number expected for limit")
	}

	page := c.QueryInt("page", 1)
	if page < 0 {
		return domain.Filter{}, fmt.Errorf("positive number expected for page")
	}

	direction := strings.ToLower(c.Query("direction"))
	var d domain.Direction

	if direction == "asc" {
		d = domain.ASC
	}

	if direction == "desc" {
		d = domain.DESC
	}

	return domain.Filter{
		Limit:     limit,
		Page:      page,
		Sort:      c.Query("sort"),
		Direction: d,
	}, nil
}
