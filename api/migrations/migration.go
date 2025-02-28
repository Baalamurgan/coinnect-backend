package migrations

import (
	"github.com/Baalamurgan/coin-selling-backend/api/db"
	"github.com/Baalamurgan/coin-selling-backend/pkg/models"
)

func Migrate() {
	database := db.GetDB()
	database.Raw("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	database.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Item{},
		&models.Detail{},
		&models.Orders{},
		&models.OrderItem{},
		&models.ShippingDetails{},
	)
}
