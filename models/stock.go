package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StockTransaction struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey"`
	ProductID       uuid.UUID
	SubVariantID    uuid.UUID
	Quantity        decimal.Decimal
	TransactionType string
	TransactionDate time.Time
}
