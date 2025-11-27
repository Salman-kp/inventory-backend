package services

import (
	"errors"
	"time"

	"inventory-backend/config"
	"inventory-backend/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockChangeRequest struct {
	ProductID    string `json:"product_id" binding:"required"`     // Product.ID (UUID)
	SubVariantID string `json:"sub_variant_id" binding:"required"` // SubVariant.ID (UUID)
	Quantity     string `json:"quantity" binding:"required"`       // decimal as string
}

type StockReport struct {
	Transactions []models.StockTransaction `json:"transactions"`
	TotalIn      decimal.Decimal           `json:"total_in"`
	TotalOut     decimal.Decimal           `json:"total_out"`
	Net          decimal.Decimal           `json:"net"`
}

func AddStock(req StockChangeRequest) (*models.StockTransaction, error) {
	return changeStock(req, "IN")
}

func RemoveStock(req StockChangeRequest) (*models.StockTransaction, error) {
	return changeStock(req, "OUT")
}

func changeStock(req StockChangeRequest, txType string) (*models.StockTransaction, error) {
	db := config.DB

	productUUID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, errors.New("invalid product_id")
	}
	subVariantUUID, err := uuid.Parse(req.SubVariantID)
	if err != nil {
		return nil, errors.New("invalid sub_variant_id")
	}
	qty, err := decimal.NewFromString(req.Quantity)
	if err != nil {
		return nil, errors.New("invalid quantity")
	}
	if qty.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("quantity must be positive")
	}

	var resultTx models.StockTransaction

	err = db.Transaction(func(tx *gorm.DB) error {
		// Lock sub-variant row
		var subVariant models.SubVariant
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&subVariant, "id = ? AND product_id = ?", subVariantUUID, productUUID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("sub_variant not found for product")
			}
			return err
		}

		// Lock product row
		var product models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&product, "id = ?", productUUID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("product not found")
			}
			return err
		}

		// Update stock
		if txType == "IN" {
			subVariant.Stock = subVariant.Stock.Add(qty)
			product.TotalStock = product.TotalStock.Add(qty)
		} else {
			if subVariant.Stock.LessThan(qty) {
				return errors.New("insufficient stock; cannot make stock negative")
			}
			subVariant.Stock = subVariant.Stock.Sub(qty)
			product.TotalStock = product.TotalStock.Sub(qty)
		}

		if err := tx.Save(&subVariant).Error; err != nil {
			return err
		}
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// Save transaction (store positive quantity; use type to distinguish)
		resultTx = models.StockTransaction{
			ID:              uuid.New(),
			ProductID:       product.ID,
			SubVariantID:    subVariant.ID,
			Quantity:        qty,
			TransactionType: txType,
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&resultTx).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resultTx, nil
}

func GetStockReport(from, to time.Time, page, limit int) (*StockReport, error) {
	db := config.DB

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	var txs []models.StockTransaction
	err := db.
		Where("transaction_date BETWEEN ? AND ?", from, to).
		Order("transaction_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error
	if err != nil {
		return nil, err
	}

	totalIn := decimal.Zero
	totalOut := decimal.Zero

	for _, t := range txs {
		if t.TransactionType == "IN" {
			totalIn = totalIn.Add(t.Quantity)
		} else if t.TransactionType == "OUT" {
			totalOut = totalOut.Add(t.Quantity)
		}
	}

	report := &StockReport{
		Transactions: txs,
		TotalIn:      totalIn,
		TotalOut:     totalOut,
		Net:          totalIn.Sub(totalOut),
	}

	return report, nil
}
