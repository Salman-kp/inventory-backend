package migrations

import (
	"inventory-backend/config"
	"inventory-backend/models"
	"log"
)

func AutoMigrate() {
	err := config.DB.AutoMigrate(
		&models.Product{},
		&models.Variant{},
		&models.VariantOption{},
		&models.SubVariant{},
		&models.StockTransaction{},
	)
	if err != nil {
		log.Fatal("AutoMigrate error:", err)
	}

	log.Println("✅ Auto migration completed")
}
