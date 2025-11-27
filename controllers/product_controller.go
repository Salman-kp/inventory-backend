package controllers

import (
	"net/http"
	"strconv"

	"inventory-backend/models"
	"inventory-backend/services"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service *services.ProductService
}

func NewProductController(s *services.ProductService) *ProductController {
	return &ProductController{service: s}
}

// POST /products
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var body models.Product

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	if err := pc.service.CreateProduct(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create product: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, body)
}

// GET /products?page=1&limit=10
func (pc *ProductController) ListProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	products, err := pc.service.ListProducts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list products: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"page":  page,
		"limit": limit,
	})
}