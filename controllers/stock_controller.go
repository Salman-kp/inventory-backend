package controllers

import (
	"net/http"
	"strconv"
	"time"

	"inventory-backend/services"

	"github.com/gin-gonic/gin"
)

// POST /api/stock/in
func AddStockHandler(c *gin.Context) {
	var req services.StockChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid request body",
			Details: err.Error(),
		})
		return
	}

	tx, err := services.AddStock(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "failed to add stock",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "stock added successfully",
		"data":    tx,
	})
}

// POST /api/stock/out
func RemoveStockHandler(c *gin.Context) {
	var req services.StockChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid request body",
			Details: err.Error(),
		})
		return
	}

	tx, err := services.RemoveStock(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "failed to remove stock",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "stock removed successfully",
		"data":    tx,
	})
}

// GET /api/stock/report?from=2025-11-01&to=2025-11-30&page=1&limit=20
func StockReportHandler(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "from and to query params are required (YYYY-MM-DD)",
		})
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid from date format",
			Details: err.Error(),
		})
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid to date format",
			Details: err.Error(),
		})
		return
	}
	to = to.Add(24 * time.Hour) // include full 'to' day

	page := 1
	limit := 20
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			page = p
		}
	}
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil {
			limit = l
		}
	}

	report, err := services.GetStockReport(from, to, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "failed to get stock report",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": report.Transactions,
		"total_in":     report.TotalIn,
		"total_out":    report.TotalOut,
		"net":          report.Net,
		"page":         page,
		"limit":        limit,
	})
}
