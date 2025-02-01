package main

import (
	"log"
	"strings"

	"github.com/Baalamurgan/coin-selling-backend/api/db"
	"github.com/Baalamurgan/coin-selling-backend/api/migrations"
	"github.com/Baalamurgan/coin-selling-backend/api/routes"
	"github.com/Baalamurgan/coin-selling-backend/api/utils"
	"github.com/Baalamurgan/coin-selling-backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func main() {
	// tracer.Start()
	// defer tracer.Stop()

	// Import environment variables
	utils.ImportEnv()

	// Load the configuration
	config.LoadConfig()

	// Init Validators
	utils.InitValidators()

	// Init Redis
	// utils.InitRedis()

	// Create a new Fiber app instance
	app := fiber.New(fiber.Config{})

	// Register routes
	app.Get("/", healthCheck)
	app.Get("/health", healthCheck)

	app.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api")
		}}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	if config.MIGRATE {
		migrations.Migrate()
	}

	// Initialize the database connection
	db.GetDB()

	routes.SetupRoutes(app)

	// Start the server on port 8080
	log.Println("Server started on http://localhost:8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
