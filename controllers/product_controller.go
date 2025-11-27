package controllers

import (
	"net/http"
	"strconv"

	"inventory-backend/services"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// POST /api/products
func CreateProductHandler(c *gin.Context) {
	var req services.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid request body",
			Details: err.Error(),
		})
		return
	}

	product, err := services.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "failed to create product",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "product created successfully",
		"data":    product,
	})
}

// GET /api/products?page=1&limit=10
func ListProductsHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	products, total, err := services.ListProducts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "failed to list products",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       products,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}
