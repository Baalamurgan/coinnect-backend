package orders

import (
	"errors"
	"fmt"

	"github.com/Baalamurgan/coin-selling-backend/api/db"
	"github.com/Baalamurgan/coin-selling-backend/api/schemas"
	"github.com/Baalamurgan/coin-selling-backend/api/utils"
	"github.com/Baalamurgan/coin-selling-backend/api/views"
	"github.com/Baalamurgan/coin-selling-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	user_id, err := uuid.Parse(req.UserID)
	if err != nil {
		return views.BadRequest(c)
	}

	var user models.User
	if err := db.GetDB().Table("users").Where("id = ?", user_id).First(&user).Error; err != nil {
		return views.BadRequestWithMessage(c, "user does not exist")
	}

	newOrder := models.Orders{
		UserID: user_id,
	}
	if err := db.GetDB().Create(&newOrder).Error; err != nil {
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
		return err
	}

	orderItem := models.OrderItem{
		OrderID:            order_id,
		ItemID:             item_id,
		BillableAmount:     item.Price,
		BillableAmountPaid: 0,
		Quantity:           1,
		OrderItemStatus:    "pending",
	}

	if err := db.GetDB().Model(&models.OrderItem{}).Create(&orderItem).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	return views.StatusOK(c, orderItem)
}

func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.GetDB().Where("id = ?", id).Delete(&models.Orders{}).Error; err != nil {
		return views.InternalServerError(c, err)
	}
	return views.StatusOK(c, "order deleted")
}
