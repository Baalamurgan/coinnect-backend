package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Baalamurgan/coin-selling-backend/api/db"
	"github.com/Baalamurgan/coin-selling-backend/api/schemas"
	"github.com/Baalamurgan/coin-selling-backend/api/utils"
	"github.com/Baalamurgan/coin-selling-backend/api/views"
	"github.com/Baalamurgan/coin-selling-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
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

	nameQuery := c.Query("name", "")
	emailQuery := c.Query("email", "")
	status := c.Query("status", "")
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

	if nameQuery != "" {
		var userIDs []uuid.UUID
		var searchMatchedUsers []models.User
		if err := db.GetDB().Model(&models.User{}).Where("username ILIKE ?", "%"+nameQuery+"%").Find(&searchMatchedUsers).Error; err == nil {
			for _, searchMatchedUser := range searchMatchedUsers {
				userIDs = append(userIDs, searchMatchedUser.ID)
			}
			dbQuery = dbQuery.Where("user_id IN ?", userIDs)
		}
	}

	if emailQuery != "" {
		var userIDs []uuid.UUID
		var emailMatchedUsers []models.User
		if err := db.GetDB().Model(&models.User{}).Where("email = ?", emailQuery).Find(&emailMatchedUsers).Error; err == nil {
			for _, emailMatchedUser := range emailMatchedUsers {
				userIDs = append(userIDs, emailMatchedUser.ID)
			}
			dbQuery = dbQuery.Where("user_id IN ?", userIDs)
		}
	}

	if status != "" {
		statuses := strings.Split(status, ",")
		if len(statuses) > 0 {
			dbQuery = dbQuery.Where("status IN ?", statuses)
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
	order_id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return views.BadRequest(c)
	}

	if err := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Preload("OrderItems").First(&order).Error; err != nil {
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
		fmt.Println(c) // no throwing error because schema is optional
	}

	orderDBQuery := db.GetDB().Model(&models.Orders{})
	newOrder := models.Orders{}

	if req.UserID != "" {
		user_id, err := uuid.Parse(req.UserID)
		if err != nil {
			return views.BadRequest(c)
		}

		if err := db.GetDB().Table("users").Where("id = ?", user_id).First(&models.User{}).Error; err != nil {
			return views.BadRequestWithMessage(c, "user does not exist")
		}

		// EXPLAIN: commented since when admin confirms an order, user side fetch still gives old order_id
		// var existingOrder models.Orders
		// if err := orderDBQuery.Where("user_id = ?", user_id).First(&existingOrder).Error; err == nil {
		// 	return views.StatusOK(c, existingOrder)
		// }

		newOrder.UserID = user_id
	}

	if err := orderDBQuery.Create(&newOrder).Error; err != nil {
		return views.InternalServerError(c, err)
	}
	return views.ObjectCreated(c, newOrder)

}

func DeleteOrder(c *fiber.Ctx) error {
	order_id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return views.BadRequest(c)
	}

	result := db.GetDB().Where("id = ?", order_id).Delete(&models.Orders{})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	return views.StatusOK(c, "order deleted")
}

func AddItemToOrder(c *fiber.Ctx) error {
	var req schemas.AddItemToOrder
	if err := c.BodyParser(&req); err != nil {
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

	var order models.Orders
	if err := db.GetDB().
		Model(&models.Orders{}).
		Where("id = ?", order_id).First(&order).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	if strings.Compare(order.Status, "pending") != 0 {
		return views.BadRequestWithMessage(c, "order confirmed already")
	}

	itemBillableAmount := item.Price*float64(quantity) + item.Price*float64(quantity)*float64((item.GST/100))

	metadata, _ := json.Marshal(map[string]interface{}{
		"category_id": item.CategoryID,
		"name":        item.Name,
		"description": item.Description,
		"year":        item.Year,
		"sku":         item.SKU,
		"image_url":   item.ImageURL,
		"stock":       item.Stock,
		"sold":        item.Sold,
		"price":       item.Price,
		"gst":         item.GST,
		"details":     item.Details,
	})

	orderItem := models.OrderItem{
		OrderID:            order_id,
		ItemID:             item_id,
		BillableAmount:     itemBillableAmount,
		BillableAmountPaid: 0,
		Quantity:           quantity,
		OrderItemStatus:    "pending",
		MetaData:           datatypes.JSON(metadata),
	}

	if err := db.GetDB().Model(&models.OrderItem{}).Create(&orderItem).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	result := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Updates(map[string]interface{}{
		"billable_amount": order.BillableAmount + itemBillableAmount,
	})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	itemResult := db.GetDB().Model(&models.Item{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"sold":  gorm.Expr("sold + ?", orderItem.Quantity),
		"stock": gorm.Expr("stock - ?", orderItem.Quantity),
	})
	if itemResult.Error != nil {
		log.Println(itemResult.Error)
	} else if itemResult.RowsAffected == 0 {
		log.Println("Stock update failed for item:", orderItem.ItemID)
	}

	return views.StatusOK(c, orderItem)
}

func DeleteOrderItemFromOrder(c *fiber.Ctx) error {
	order_id, err := uuid.Parse(c.Params("order_id"))
	if err != nil {
		return views.BadRequest(c)
	}
	order_item_id, err := uuid.Parse(c.Params("order_item_id"))
	if err != nil {
		return views.BadRequest(c)
	}

	orderDBQuery := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id)

	var order models.Orders
	if err := orderDBQuery.First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	if strings.Compare(order.Status, "pending") != 0 {
		return views.BadRequestWithMessage(c, "order confirmed already")
	}

	var orderItem models.OrderItem
	if err := db.GetDB().Model(&models.OrderItem{}).First(&orderItem, order_item_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	if err := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Updates(map[string]interface{}{
		"billable_amount": order.BillableAmount - orderItem.BillableAmount,
	}).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	result := db.GetDB().Where("order_id = ? AND id = ?", order_id, order_item_id).Delete(&models.OrderItem{})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	itemResult := db.GetDB().Model(&models.Item{}).Where("id = ?", orderItem.ItemID).Updates(map[string]interface{}{
		"sold":  gorm.Expr("sold - ?", orderItem.Quantity),
		"stock": gorm.Expr("stock + ?", orderItem.Quantity),
	})
	if itemResult.Error != nil {
		log.Println(itemResult.Error)
	} else if itemResult.RowsAffected == 0 {
		log.Println("Stock update failed for item:", orderItem.ItemID)
	}

	return views.StatusOK(c, "order item deleted")
}

func ConfirmOrder(c *fiber.Ctx) error {
	var req schemas.ConfirmOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	order_id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return views.BadRequest(c)
	}

	user_id, err := uuid.Parse(req.UserID)
	if err != nil {
		return views.BadRequest(c)
	}

	if err := db.GetDB().Where("id = ?", user_id).First(&models.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	var order models.Orders
	if err := db.GetDB().Where("id = ?", order_id).Preload("OrderItems").Find(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	if strings.Compare(order.Status, "cancelled") == 0 {
		return views.BadRequestWithMessage(c, "order has already been cancelled")
	}

	if len(order.OrderItems) <= 0 {
		return views.BadRequestWithMessage(c, "order invalid")
	}

	result := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Updates(map[string]interface{}{
		"status":      "booked",
		"user_id":     user_id,
		"status_date": time.Now().Unix(),
	})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	return views.StatusOK(c, "order confirmed")
}

func MarkOrderAsPaid(c *fiber.Ctx) error {
	var req schemas.MarkOrderAsPaidRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	order_id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return views.BadRequest(c)
	}

	user_id, err := uuid.Parse(req.UserID)
	if err != nil {
		return views.BadRequest(c)
	}

	if err := db.GetDB().Where("id = ?", user_id).First(&models.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	var order models.Orders
	if err := db.GetDB().Where("id = ?", order_id).Find(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	if strings.Compare(order.Status, "cancelled") == 0 {
		return views.BadRequestWithMessage(c, "order has already been cancelled")
	}

	if strings.Compare(order.Status, "booked") != 0 {
		return views.BadRequestWithMessage(c, "order must be confirmed before payment")
	}

	if req.BillableAmountPaid < 0 {
		return views.BadRequestWithMessage(c, "invalid payment amount")
	}

	result := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Updates(map[string]interface{}{
		"status":               "paid",
		"billable_amount_paid": req.BillableAmountPaid,
		"status_date":          time.Now().Unix(),
	})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	return views.StatusOK(c, "order paid")
}

func MarkOrderAsShipped(c *fiber.Ctx) error {
	var req schemas.MarkOrderAsShippedRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	order_id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return views.BadRequest(c)
	}

	user_id, err := uuid.Parse(req.UserID)
	if err != nil {
		return views.BadRequest(c)
	}

	if err := db.GetDB().Where("id = ?", user_id).First(&models.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	var order models.Orders
	if err := db.GetDB().Where("id = ?", order_id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	if strings.Compare(order.Status, "cancelled") == 0 {
		return views.BadRequestWithMessage(c, "order has already been cancelled")
	}

	if strings.Compare(order.Status, "paid") != 0 {
		return views.BadRequestWithMessage(c, "order has not been paid yet")
	}

	var shippingDetails models.ShippingDetails
	shippingDetails.OrderID = order_id
	shippingDetails.UserID = user_id
	shippingDetails.ShippingName = req.ShippingName
	shippingDetails.ShippingID = req.ShippingID
	shippingDetails.ShippingDate = req.ShippingDate

	if err := db.GetDB().Create(&shippingDetails).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	result := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Updates(map[string]interface{}{
		"status":      "shipped",
		"shipping_id": shippingDetails.ID,
		"status_date": time.Now().Unix(),
	})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	return views.StatusOK(c, "order shipped")
}

func CancelOrder(c *fiber.Ctx) error {
	var req schemas.CancelOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	order_id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return views.BadRequest(c)
	}

	user_id, err := uuid.Parse(req.UserID)
	if err != nil {
		return views.BadRequest(c)
	}

	if err := db.GetDB().Where("id = ?", user_id).First(&models.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	var order models.Orders
	if err := db.GetDB().Where("id = ?", order_id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	if strings.Compare(order.Status, "cancelled") == 0 {
		return views.BadRequestWithMessage(c, "order has already been cancelled")
	}

	if strings.Compare(order.Status, "delivered") == 0 {
		return views.BadRequestWithMessage(c, "delivered orders cannot be cancelled")
	}

	result := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Updates(map[string]interface{}{
		"status":              "cancelled",
		"cancellation_reason": req.CancellationReason,
		"status_date":         time.Now().Unix(),
	})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	return views.StatusOK(c, "order cancelled")
}

func MarkOrderAsDelivered(c *fiber.Ctx) error {
	var req schemas.MarkOrderAsDeliveredRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	order_id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return views.BadRequest(c)
	}

	user_id, err := uuid.Parse(req.UserID)
	if err != nil {
		return views.BadRequest(c)
	}

	if err := db.GetDB().Where("id = ?", user_id).First(&models.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	var order models.Orders
	if err := db.GetDB().Where("id = ?", order_id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	if strings.Compare(order.Status, "cancelled") == 0 {
		return views.BadRequestWithMessage(c, "order has already been cancelled")
	}

	if strings.Compare(order.Status, "shipped") != 0 {
		return views.BadRequestWithMessage(c, "order has not been shipped yet")
	}

	var deliveryDetails models.DeliveryDetails
	deliveryDetails.OrderID = order_id
	deliveryDetails.UserID = user_id
	deliveryDetails.DeliveryPersonName = req.DeliveryPersonName
	deliveryDetails.DeliveryID = req.DeliveryID
	deliveryDetails.DeliveryDate = req.DeliveryDate

	if err := db.GetDB().Create(&deliveryDetails).Error; err != nil {
		return views.InternalServerError(c, err)
	}

	result := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Updates(map[string]interface{}{
		"status":      "delivered",
		"delivery_id": deliveryDetails.ID,
		"status_date": time.Now().Unix(),
	})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	return views.StatusOK(c, "order delivered")
}

func RestoreOrder(c *fiber.Ctx) error {
	var req schemas.RestoreOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return views.InvalidParams(c)
	}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return views.InvalidParams(c)
	}

	order_id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return views.BadRequest(c)
	}

	user_id, err := uuid.Parse(req.UserID)
	if err != nil {
		return views.BadRequest(c)
	}

	if err := db.GetDB().Where("id = ?", user_id).First(&models.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	var order models.Orders
	if err := db.GetDB().Where("id = ?", order_id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.RecordNotFound(c)
		}
		return views.InternalServerError(c, err)
	}

	if strings.Compare(order.Status, "cancelled") != 0 {
		return views.BadRequestWithMessage(c, "order is not cancelled")
	}

	result := db.GetDB().Model(&models.Orders{}).Where("id = ?", order_id).Updates(map[string]interface{}{
		"status":      "pending",
		"status_date": time.Now().Unix(),
	})
	if result.Error != nil {
		return views.InternalServerError(c, result.Error)
	} else if result.RowsAffected == 0 {
		return views.RecordNotFound(c)
	}

	return views.StatusOK(c, "order restored")
}
