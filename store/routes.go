package store

import "github.com/gofiber/fiber/v2"

func Routes(f *fiber.App, s Service) {
	f.Post("/v1/products", addProduct(s))
	f.Get("/v1/products", listProducts(s))
	f.Get("/v1/products/:id", findProduct(s))
	f.Put("/v1/products/:id", updateProduct(s))
	f.Delete("/v1/products/:id", deleteProduct(s))
}
