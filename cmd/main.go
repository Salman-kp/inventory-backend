package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"inventory-backend/config"
	"inventory-backend/migrations"
	"inventory-backend/routes"
)

func main() {
	config.ConnectDatabase()
	migrations.AutoMigrate()

	router := gin.Default()
	routes.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Server running at http://localhost%s", addr)
	router.Run(addr)
}
