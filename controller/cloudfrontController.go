package controller

import (
	cdn "go-sdk/service/cdn"

	"github.com/gofiber/fiber/v2"
)

func CloudfrontController(router fiber.Router) {
	// GET 요청으로 아래와 같이 요청이 들어오면 될 듯
	// localhost:3000/api/v1/securityGroup?name=search&type=alb
	// query parameter 의 type 을 바탕으로 switch 문 작성
	router.Post("/", cdn.CreateInvalidation)
}
