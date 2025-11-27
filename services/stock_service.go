package services

import (
	"errors"
	"time"

	"inventory-backend/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockService struct {
	db *gorm.DB
}

func NewStockService(db *gorm.DB) *StockService {
	return &StockService{db: db}
}

// AddStock performs a stock-in operation with row-level locking
func (s *StockService) AddStock(productID, subVariantID uuid.UUID, qty decimal.Decimal) error {
	if qty.LessThanOrEqual(decimal.Zero) {
		return errors.New("quantity must be positive")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var sub models.SubVariant

		// Row-level lock the sub_variant row
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", subVariantID).
			First(&sub).Error; err != nil {
			return err
		}

		newStock := sub.Stock.Add(qty)

		stockTx := &models.StockTransaction{
			ID:              uuid.New(),
			ProductID:       productID,
			SubVariantID:    subVariantID,
			Quantity:        qty,
			TransactionType: "IN",
			TransactionDate: time.Now(),
		}

		if err := tx.Create(stockTx).Error; err != nil {
			return err
		}

		// Update SubVariant stock
		if err := tx.Model(&models.SubVariant{}).
			Where("id = ?", subVariantID).
			Update("stock", newStock).Error; err != nil {
			return err
		}

		// Optional: update Product.TotalStock as sum
		if err := tx.Model(&models.Product{}).
			Where("id = ?", productID).
			Update("total_stock", gorm.Expr("total_stock + ?", qty)).Error; err != nil {
			return err
		}

		return nil
	})
}

// RemoveStock performs a stock-out operation with negative stock prevention
func (s *StockService) RemoveStock(productID, subVariantID uuid.UUID, qty decimal.Decimal) error {
	if qty.LessThanOrEqual(decimal.Zero) {
		return errors.New("quantity must be positive")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var sub models.SubVariant

		// Row-level lock
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", subVariantID).
			First(&sub).Error; err != nil {
			return err
		}

		if sub.Stock.LessThan(qty) {
			return errors.New("insufficient stock, cannot go negative")
		}

		newStock := sub.Stock.Sub(qty)

		stockTx := &models.StockTransaction{
			ID:              uuid.New(),
			ProductID:       productID,
			SubVariantID:    subVariantID,
			Quantity:        qty.Neg(),
			TransactionType: "OUT",
			TransactionDate: time.Now(),
		}

		if err := tx.Create(stockTx).Error; err != nil {
			return err
		}

		// Update SubVariant stock
		if err := tx.Model(&models.SubVariant{}).
			Where("id = ?", subVariantID).
			Update("stock", newStock).Error; err != nil {
			return err
		}

		// Update Product.TotalStock
		if err := tx.Model(&models.Product{}).
			Where("id = ?", productID).
			Update("total_stock", gorm.Expr("total_stock - ?", qty)).Error; err != nil {
			return err
		}

		return nil
	})
}

// StockReport returns all stock transactions between two dates (inclusive)
func (s *StockService) StockReport(from, to string) ([]models.StockTransaction, error) {
	var logs []models.StockTransaction

	err := s.db.
		Where("transaction_date BETWEEN ? AND ?", from, to).
		Order("transaction_date DESC").
		Find(&logs).Error

	return logs, err
}
