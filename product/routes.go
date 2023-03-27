package product

import "github.com/gofiber/fiber/v2"

func Routes(f *fiber.App, s Service) {
	f.Post("/v1/products", addProduct(s))
}
