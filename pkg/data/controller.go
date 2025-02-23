package data

import (
	"encoding/json"
	"io"

	"github.com/Baalamurgan/coin-selling-backend/api/db"
	"github.com/Baalamurgan/coin-selling-backend/api/views"
	"github.com/Baalamurgan/coin-selling-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func Populate(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return views.BadRequestWithMessage(c, "Failed to get file")
	}

	fileContent, err := file.Open()
	if err != nil {
		return views.BadRequestWithMessage(c, "Failed to open file")
	}
	defer fileContent.Close()

	data, err := io.ReadAll(fileContent)
	if err != nil {
		return views.BadRequestWithMessage(c, "Failed to read file")
	}

	var categories []models.Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return views.BadRequestWithMessage(c, "Invalid JSON format")
	}

	if err := db.GetDB().Model(&models.Category{}).Create(&categories).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	return views.StatusOK(c, "Data inserted successfully")
}
