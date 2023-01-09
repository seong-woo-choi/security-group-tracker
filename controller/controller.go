package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func MainController(app *fiber.App) {
	apiV1Group := app.Group("/api/v1", logger.New())
	securityGroup := apiV1Group.Group("/securityGroup")

	SecurityGroupController(securityGroup)
}
