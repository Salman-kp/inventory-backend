package routes

import (
	"inventory-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Products
		api.POST("/products", controllers.CreateProductHandler)
		api.GET("/products", controllers.ListProductsHandler)

		// Stock
		api.POST("/stock/in", controllers.AddStockHandler)
		api.POST("/stock/out", controllers.RemoveStockHandler)

		// Stock report
		api.GET("/stock/report", controllers.StockReportHandler)
	}
}
