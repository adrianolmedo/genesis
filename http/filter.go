package http

import (
	"fmt"
	"strconv"

	domain "github.com/adrianolmedo/genesis"

	"github.com/gofiber/fiber/v2"
)

func getFilter(c *fiber.Ctx) (domain.Filter, error) {
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if limit < 0 || err != nil {
		return domain.Filter{}, fmt.Errorf("positive number expected for limit")
	}

	page, err := strconv.Atoi(c.Query("page", "0"))
	if page < 0 || err != nil {
		return domain.Filter{}, fmt.Errorf("positive number expected for page")
	}

	return domain.Filter{
		Limit:     limit,
		Page:      page,
		Sort:      c.Query("sort"),
		Direction: c.Query("direction"),
	}, nil
}
