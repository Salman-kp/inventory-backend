package main

import (
	"fmt"
	"log"
	"os"

	"inventory-backend/config"
	"inventory-backend/controllers"
	"inventory-backend/migrations"
	"inventory-backend/routes"
	"inventory-backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	migrations.AutoMigrate()

	db := config.DB

	productService := services.NewProductService(db)
	stockService := services.NewStockService(db)

	productController := controllers.NewProductController(productService)
	stockController := controllers.NewStockController(stockService)

	router := gin.Default()
	routes.RegisterRoutes(router, productController, stockController)
	
	
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Server running at http://localhost%s", addr)
	router.Run(addr)
}
