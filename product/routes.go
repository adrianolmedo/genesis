package product

import "github.com/gofiber/fiber/v2"

func Routes(f *fiber.App) {
	f.Get("/v1/products", addProduct())
}
