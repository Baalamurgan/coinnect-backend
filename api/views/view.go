package views

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func InvalidParams(c *fiber.Ctx) error {
	return c.
		Status(400).
		JSON(fiber.Map{
			"message": "invalid params",
		})
}

func InternalServerError(c *fiber.Ctx, err error) error {
	log.Println(err)
	return c.
		Status(500).
		JSON(fiber.Map{
			"status":  "fail",
			"message": "something went wrong",
		})
}

func StatusOK(c *fiber.Ctx, data interface{}) error {
	return c.
		Status(200).
		JSON(fiber.Map{
			"status": "true",
			"data":   data,
		})

}

func ObjectCreated(c *fiber.Ctx, data interface{}) error {
	return c.
		Status(201).
		JSON(fiber.Map{
			"status": "true",
			"data":   data,
		})
}

func RecordNotFound(c *fiber.Ctx) error {
	return c.
		Status(404).
		JSON(fiber.Map{
			"status":  "fail",
			"message": "not found",
		})
}

func UnAuthorisedView(c *fiber.Ctx) error {
	return c.
		Status(401).
		JSON(fiber.Map{
			"status":  "error",
			"message": "un authorized",
		})
}

func UnAuthorisedViewWithMessage(c *fiber.Ctx, message string) error {
	return c.
		Status(401).
		JSON(fiber.Map{
			"status":  "error",
			"message": message})
}

func ForbiddenView(c *fiber.Ctx) error {
	return c.
		Status(403).
		JSON(fiber.Map{
			"status":  "fail",
			"message": "forbidden",
		})
}

func BadRequest(c *fiber.Ctx) error {
	return c.
		Status(400).JSON(fiber.Map{
		"status":  "fail",
		"message": "bad request",
	})
}

func BadRequestWithMessage(c *fiber.Ctx, message string) error {
	return c.
		Status(400).JSON(fiber.Map{
		"status":  "error",
		"message": message})
}

func Conflict(c *fiber.Ctx) error {
	return c.
		Status(409).JSON(fiber.Map{
		"status":  "fail",
		"message": "conflict",
	})
}

func ConflictWithMessage(c *fiber.Ctx, message string) error {
	return c.
		Status(409).JSON(fiber.Map{
		"status":  "fail",
		"message": message,
	})
}
