package routes

import (
	"github.com/Baalamurgan/coin-selling-backend/pkg/auth"
	"github.com/Baalamurgan/coin-selling-backend/pkg/category"
	"github.com/Baalamurgan/coin-selling-backend/pkg/item"
	"github.com/Baalamurgan/coin-selling-backend/pkg/orders"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	//Auth
	authGroup := v1.Group("/auth")
	authGroup.Post("/signup", auth.Signup)
	authGroup.Post("/login", auth.Login)
	authGroup.Post("/profile", auth.GetUser)
	authGroup.Get("/users", auth.GetAllUsers)
	authGroup.Post("/update/:id", auth.UpdateUser)

	// Category
	categoryGroup := v1.Group("/category")
	categoryGroup.Get("/", category.GetAllCategories)
	categoryGroup.Get("/:id", category.GetCategoryByID)
	categoryGroup.Get("/:id/all", category.GetAllCategoriesByParentCategoryID)
	categoryGroup.Post("/", category.CreateCategory)
	categoryGroup.Put("/:id", category.UpdateCategory)
	categoryGroup.Delete("/:id", category.DeleteCategory)

	// Sub Category
	// subCategoryGroup := v1.Group("/sub_category")
	// subCategoryGroup.Get("/", sub_category.GetAllSubCategories)
	// subCategoryGroup.Get("/category/:category_id", sub_category.GetAllSubCategoriesByCategoryID)
	// subCategoryGroup.Get("/:id", sub_category.GetSubCategoryByID)
	// subCategoryGroup.Post("/:category_id", sub_category.CreateSubCategory)
	// subCategoryGroup.Put("/:id", sub_category.UpdateSubCategory)
	// subCategoryGroup.Delete("/:id", sub_category.DeleteSubCategory)

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
	orderGroup.Post("/item/add", orders.AddItemToOrder)
	orderGroup.Delete("/:id", orders.DeleteOrder)

	// coinGroup := itemGroup.Group("/coin")
}
