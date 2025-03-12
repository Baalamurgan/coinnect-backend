package migrations

import (
	"log"

	"github.com/Baalamurgan/coin-selling-backend/api/db"
	"github.com/Baalamurgan/coin-selling-backend/pkg/models"
)

func Migrate() {
	database := db.GetDB()
	err := database.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatalf("Error enabling UUID extension: %v", err)
	}
	database.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Item{},
		&models.Detail{},
		&models.Orders{},
		&models.OrderItem{},
		&models.ShippingDetails{},
		&models.DeliveryDetails{},
	)
}
