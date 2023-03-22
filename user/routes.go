package user

import "github.com/gofiber/fiber/v2"

func Routes(f *fiber.App, s Service) {
	f.Get("/v1/users/:id", findUser(s))
}
