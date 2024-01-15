package http

import (
	domain "github.com/adrianolmedo/genesis"

	"github.com/gofiber/fiber/v2"
)

func getFilter(c *fiber.Ctx) (*domain.Filter, error) {
	f := domain.NewFilter()

	err := f.SetLimit(c.QueryInt("limit"))
	if err != nil {
		return nil, err
	}

	err = f.SetPage(c.QueryInt("page"))
	if err != nil {
		return nil, err
	}

	f.SetSort(c.Query("sort", "created_at"))
	f.SetDirection(c.Query("direction"))
	return f, nil
}
