package controllers

import (
	"net/http"

	"inventory-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StockController struct {
	service *services.StockService
}

func NewStockController(s *services.StockService) *StockController {
	return &StockController{service: s}
}

type StockRequest struct {
	ProductID    string `json:"product_id" binding:"required"`
	SubVariantID string `json:"sub_variant_id" binding:"required"`
	Quantity     string `json:"quantity" binding:"required"`
}

// POST /stock/add
func (sc *StockController) AddStock(c *gin.Context) {
	var req StockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	productUUID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	subUUID, err := uuid.Parse(req.SubVariantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sub_variant_id"})
		return
	}

	qty, err := decimal.NewFromString(req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
		return
	}

	if err := sc.service.AddStock(productUUID, subUUID, qty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock added successfully"})
}

// POST /stock/remove
func (sc *StockController) RemoveStock(c *gin.Context) {
	var req StockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	productUUID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	subUUID, err := uuid.Parse(req.SubVariantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sub_variant_id"})
		return
	}

	qty, err := decimal.NewFromString(req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
		return
	}

	if err := sc.service.RemoveStock(productUUID, subUUID, qty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock removed successfully"})
}

// GET /stock/report?from=2024-01-01&to=2025-01-01
func (sc *StockController) StockReport(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from and to query params are required",
		})
		return
	}

	logs, err := sc.service.StockReport(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch stock report: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
	})
}