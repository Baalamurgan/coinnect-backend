package orders

import (
	"errors"
	"fmt"
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

func GetAllOrders(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return views.BadRequest(c)
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 0 {
		return views.BadRequest(c)
	}

	searchQuery := c.Query("search", "")
	categoryIDs := c.Query("category_ids", "")

	var parsedCategoryIDs []*uuid.UUID
	if categoryIDs != "" {
		categoryIDs := strings.Split(categoryIDs, "")
		for _, categoryID := range categoryIDs {
			parsedCategoryID, err := utils.ParseUUID(categoryID)
			if err == nil {
				parsedCategoryIDs = append(parsedCategoryIDs, parsedCategoryID)
			}
		}
	}

	dbQuery := db.GetDB().Model(&models.Orders{})

	var total int64
	var orders []models.Orders

	if searchQuery != "" {
		var searchedUser models.User
		if err := db.GetDB().Model(&models.User{}).Where("username ILIKE ? OR email ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%").First(&searchedUser).Error; err != nil {
			dbQuery = dbQuery.Where("user_id = ?", searchedUser.ID)
		}
	}

	if parsedCategoryIDs != nil {
		dbQuery = dbQuery.Where("category_id IN ?", parsedCategoryIDs)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	if err := dbQuery.Order("updated_at DESC").Scopes(utils.Paginate(page, limit)).Find(&orders).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	return views.StatusOK(c, fiber.Map{
		"orders": orders,
		"pagination": fiber.Map{
			"page":          page,
			"limit":         limit,
			"total_records": total,
			"total_pages":   utils.CalculateTotalPages(total, limit),
		},
	})
}

func GetOrderByID(c *fiber.Ctx) error {
	var order models.Orders
	id := c.Params("id")
	if err := db.GetDB().Model(&models.Orders{}).Where("id = ?", id).Preload("OrderItems").First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, order)
}

func CreateOrder(c *fiber.Ctx) error {
	var req schemas.CreateOrder
	if err := c.BodyParser(&req); err != nil {
		fmt.Println(c)
	}

	orderDBQuery := db.GetDB().Model(&models.Orders{})
	newOrder := models.Orders{}

	if req.UserID != "" {
		user_id, err := uuid.Parse(req.UserID)
		if err != nil {
			return views.BadRequest(c)
		}

		var user models.User
		if err := db.GetDB().Table("users").Where("id = ?", user_id).First(&user).Error; err != nil {
			return views.BadRequestWithMessage(c, "user does not exist")
		}

		var existingOrder models.Orders
		if err := orderDBQuery.Where("user_id = ?", user_id).First(&existingOrder).Error; err == nil {
			return views.StatusOK(c, existingOrder)
		}

		newOrder.UserID = user_id
	}

	if err := orderDBQuery.Create(&newOrder).Error; err != nil {
		return views.InternalServerError(c, err)
	}
	return views.ObjectCreated(c, newOrder)

}

func AddItemToOrder(c *fiber.Ctx) error {
	var req schemas.AddItemToOrder
	if err := c.BodyParser(&req); err != nil {
		fmt.Println(c)
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	quantity := req.Quantity
	if quantity < 1 {
		return views.BadRequest(c)
	}

	order_id, err := uuid.Parse(req.OrderID)
	if err != nil {
		return views.BadRequest(c)
	}

	item_id, err := uuid.Parse(req.ItemID)
	if err != nil {
		return views.BadRequest(c)
	}

	var item models.Item
	if err := db.GetDB().First(&item, item_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	if quantity > item.Stock {
		return views.BadRequestWithMessage(c, "requested quantity exceeds available stock")
	}

	itemBillableAmount := item.Price*float64(quantity) + item.Price*float64(quantity)*float64((item.GST/100))

	orderItem := models.OrderItem{
		OrderID:            order_id,
		ItemID:             item_id,
		BillableAmount:     itemBillableAmount,
		BillableAmountPaid: 0,
		Quantity:           quantity,
		OrderItemStatus:    "pending",
	}

	if err := db.GetDB().Model(&models.OrderItem{}).Create(&orderItem).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	var totalBillableAmount float64
	if err := db.GetDB().
		Model(&models.Orders{}).
		Where("id = ?", order_id).
		Select("billable_amount").
		Scan(&totalBillableAmount).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	if err := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Updates(map[string]interface{}{
		"billable_amount": totalBillableAmount + itemBillableAmount,
	}).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	return views.StatusOK(c, orderItem)
}

func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.GetDB().Where("id = ?", id).Delete(&models.Orders{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, "order deleted")
}
