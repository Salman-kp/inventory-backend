package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey"`
	ProductID    int64           `gorm:"unique"`
	ProductCode  string          `gorm:"unique"`
	ProductName  string
	ProductImage string
	CreatedDate  time.Time       `gorm:"autoCreateTime"`
	UpdatedDate  time.Time
	CreatedUser  uuid.UUID       `gorm:"type:uuid"`
	IsFavourite  bool
	Active       bool
	HSNCode      string
	TotalStock   decimal.Decimal `gorm:"type:numeric(20,8);default:0"`
}
