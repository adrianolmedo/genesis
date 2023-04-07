package user

import "github.com/gofiber/fiber/v2"

func Routes(f *fiber.App, s Service) {
	f.Post("v1/users", signUpUser(s))
	f.Get("v1/users", listUsers(s))
	f.Get("/v1/users/:id", findUser(s))
	f.Put("v1/users/:id", updateUser(s))
	f.Delete("v1/users/:id", deleteUser(s))
}
