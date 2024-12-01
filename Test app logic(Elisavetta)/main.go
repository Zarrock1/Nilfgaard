package main

import (
	"core_mod/controllers"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/users/:id", controllers.UserHandler)
	app.Get("/users", controllers.UsersHandler)
	app.Patch("/users/:id/name", controllers.UserUpdate)
	app.Get("/users/:id/role", controllers.UserRoles)
	// app.Patch("/users/:id/role", controller.)

	log.Println(app.Listen(":3000"))

}
