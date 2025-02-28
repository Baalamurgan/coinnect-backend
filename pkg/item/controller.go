package item

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Baalamurgan/coin-selling-backend/api/db"
	"github.com/Baalamurgan/coin-selling-backend/api/schemas"
	"github.com/Baalamurgan/coin-selling-backend/api/utils"
	"github.com/Baalamurgan/coin-selling-backend/api/views"
	"github.com/Baalamurgan/coin-selling-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetAllItems(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return views.BadRequest(c)
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 0 {
		return views.BadRequest(c)
	}

	searchQuery := c.Query("search", "")
	categoryIDs := c.Query("category_ids", "") // category_id1, category_id2, category_id3

	var parsedCategoryIDs []*uuid.UUID
	if categoryIDs != "" {
		categoryIDList := strings.Split(categoryIDs, ",")

		for _, categoryID := range categoryIDList {
			parsedCategoryID, err := utils.ParseUUID(categoryID)
			if err == nil {
				parsedCategoryIDs = append(parsedCategoryIDs, parsedCategoryID)
			}
		}
	}

	var items []models.Item
	var total int64
	dbQuery := db.GetDB().Model(&models.Item{}).Preload("Details")

	if searchQuery != "" {
		dbQuery = dbQuery.Where("name ILIKE ? OR description ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	if parsedCategoryIDs != nil {
		dbQuery = dbQuery.Where("category_id IN ?", parsedCategoryIDs)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	if err := dbQuery.Order("updated_at DESC").Scopes(utils.Paginate(page, limit)).Find(&items).Error; err != nil {
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, fiber.Map{
		"items": items,
		"pagination": fiber.Map{
			"page":          page,
			"limit":         limit,
			"total_records": total,
			"total_pages":   utils.CalculateTotalPages(total, limit),
		},
	})
}

func GetItemsByCategoryID(c *fiber.Ctx) error {
	category_id := c.Params("category_id")
	if category_id == "" {
		return views.BadRequestWithMessage(c, "category id required")
	}
	var items []models.Item

	if err := db.GetDB().Where("category_id = ?", category_id).Find(&items).Error; err != nil {
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, items)
}

func GetItemsBySubCategoryID(c *fiber.Ctx) error {
	sub_category_id := c.Params("sub_category_id")
	if sub_category_id == "" {
		return views.BadRequestWithMessage(c, "sub category id required")
	}
	var items []models.Item

	if err := db.GetDB().Where("sub_category_id = ?", sub_category_id).Find(&items).Error; err != nil {
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, items)
}

func GetItemByID(c *fiber.Ctx) error {
	var item models.Item
	id := c.Params("id")
	if err := db.GetDB().Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, item)
}

func GetItemBySlug(c *fiber.Ctx) error {
	var item models.Item
	slug := strings.ReplaceAll(c.Params("slug"), "%26", "&")

	if err := db.GetDB().Model(&models.Item{}).Preload("Details").Where("slug = ?", slug).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, item)
}

func CreateItem(c *fiber.Ctx) error {
	category_id, err := uuid.Parse(c.Params("category_id"))
	if err != nil {
		return views.BadRequest(c)
	}
	var req schemas.CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	var newItem models.Item
	newItem.CategoryID = category_id
	newItem.Name = req.Name
	newItem.Description = req.Description
	newItem.Year = req.Year
	newItem.SKU = req.SKU
	newItem.ImageURL = req.ImageURL
	newItem.Stock = req.Stock
	newItem.Sold = req.Sold
	newItem.Price = req.Price
	newItem.GST = req.GST
	newItem.Slug = utils.GenerateSlug(req.Name)

	for _, itemReq := range req.Details {
		newItem.Details = append(newItem.Details, models.Detail{
			Attribute: itemReq.Attribute,
			Value:     itemReq.Value,
		})
	}

	if err := db.GetDB().Create(&newItem).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	return views.ObjectCreated(c, newItem)
}

func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var req schemas.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	var item models.Item
	if err := db.GetDB().Table("items").Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	item.CategoryID = *req.CategoryID
	item.Name = *req.Name
	item.Description = *req.Description
	item.Year = *req.Year
	item.SKU = *req.SKU
	item.ImageURL = *req.ImageURL
	item.Stock = *req.Stock
	item.Sold = *req.Sold
	item.Price = *req.Price
	item.GST = *req.GST

	if err := db.GetDB().Save(&item).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	return views.StatusOK(c, "item updated successfully")
}

func DeleteItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return views.BadRequest(c)
	}

	result := db.GetDB().Where("id = ?", id).Delete(&models.Item{})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	return views.StatusOK(c, "item deleted")
}
