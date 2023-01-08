package main

import (
	"go-sdk/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	controller.MainController(app)
	app.Listen(":3000")
}
