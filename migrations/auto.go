package migrations

import (
	"inventory-backend/config"
	"inventory-backend/models"
)

func AutoMigrate() {
	config.DB.AutoMigrate(
		&models.Product{},
		&models.Variant{},
		&models.VariantOption{},
		&models.SubVariant{},
		&models.StockTransaction{},
	)
}
