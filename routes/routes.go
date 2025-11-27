package routes

import (
	"inventory-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	productController *controllers.ProductController,
	stockController *controllers.StockController,
) {


	// Product APIs
	r.POST("/products", productController.CreateProduct)
	r.GET("/products", productController.ListProducts)

	// Stock APIs
	r.POST("/stock/add", stockController.AddStock)
	r.POST("/stock/remove", stockController.RemoveStock)
	r.GET("/stock/report", stockController.StockReport)
}
