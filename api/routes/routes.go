package routes

import (
	"github.com/Baalamurgan/coin-selling-backend/pkg/auth"
	"github.com/Baalamurgan/coin-selling-backend/pkg/category"
	"github.com/Baalamurgan/coin-selling-backend/pkg/data"
	"github.com/Baalamurgan/coin-selling-backend/pkg/item"
	"github.com/Baalamurgan/coin-selling-backend/pkg/orders"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/populate", data.Populate)

	// Auth
	authGroup := v1.Group("/auth")
	authGroup.Post("/signup", auth.Signup)
	authGroup.Post("/login", auth.Login)
	authGroup.Get("/users", auth.GetAllUsers)
	// Profile
	profileGroup := authGroup.Group("/profile")
	profileGroup.Post("/", auth.GetUser)
	profileGroup.Post("/email", auth.GetUserByEmail)
	profileGroup.Put("/update/:id", auth.UpdateUser)

	// Category
	categoryGroup := v1.Group("/category")
	categoryGroup.Get("/", category.GetAllCategories)
	categoryGroup.Get("/:id", category.GetCategoryByID)
	categoryGroup.Get("/:id/all", category.GetAllCategoriesByParentCategoryID)
	categoryGroup.Post("/", category.CreateCategory)
	categoryGroup.Put("/:id", category.UpdateCategory)
	categoryGroup.Delete("/:id", category.DeleteCategory)

	// Item
	itemGroup := v1.Group("/item")
	itemGroup.Get("/", item.GetAllItems)
	itemGroup.Get("/category/:category_id", item.GetItemsByCategoryID)
	itemGroup.Get("/sub_category/:sub_category_id", item.GetItemsBySubCategoryID)
	itemGroup.Get("/:id", item.GetItemByID)
	itemGroup.Get("/slug/:slug", item.GetItemBySlug)
	itemGroup.Post("/:category_id", item.CreateItem)
	itemGroup.Put("/:id", item.UpdateItem)
	itemGroup.Delete("/:id", item.DeleteItem)

	// Order
	orderGroup := v1.Group("/order")
	orderGroup.Get("/", orders.GetAllOrders)
	orderGroup.Get("/:id", orders.GetOrderByID)
	orderGroup.Post("/", orders.CreateOrder)
	orderGroup.Delete("/:id", orders.DeleteOrder)
	orderGroup.Patch("/:id/edit", orders.EditOrder)
	orderGroup.Patch("/:id/confirm", orders.ConfirmOrder)
	orderGroup.Patch("/:id/cancel", orders.CancelOrder)
	orderGroup.Patch("/:id/pay", orders.MarkOrderAsPaid)
	orderGroup.Patch("/:id/ship", orders.MarkOrderAsShipped)
	orderGroup.Patch("/:id/deliver", orders.MarkOrderAsDelivered)
	orderGroup.Patch("/:id/restore", orders.RestoreOrder)
	// Order Item
	orderItemGroup := orderGroup.Group("/item")
	orderItemGroup.Post("/add", orders.AddItemToOrder)
	orderItemGroup.Delete("/:order_id/:order_item_id", orders.DeleteOrderItemFromOrder)
}
