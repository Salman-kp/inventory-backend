package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StockTransaction struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID       uuid.UUID       `gorm:"type:uuid;index" json:"product_id"`
	SubVariantID    uuid.UUID       `gorm:"type:uuid;index" json:"sub_variant_id"`
	Quantity        decimal.Decimal `gorm:"type:numeric(20,8)" json:"quantity"`
	TransactionType string          `json:"transaction_type"` // "IN" or "OUT"
	TransactionDate time.Time       `json:"transaction_date"`
}