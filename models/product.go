package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID    int64           `gorm:"unique" json:"product_id"`
	ProductCode  string          `gorm:"unique" json:"product_code"`
	ProductName  string          `json:"product_name"`
	ProductImage string          `json:"product_image"`
	CreatedDate  time.Time       `gorm:"autoCreateTime" json:"created_date"`
	UpdatedDate  time.Time       `gorm:"autoUpdateTime" json:"updated_date"`
	CreatedUser  uuid.UUID       `gorm:"type:uuid" json:"created_user"`
	IsFavourite  bool            `json:"is_favourite"`
	Active       bool            `json:"active"`
	HSNCode      string          `json:"hsn_code"`
	TotalStock   decimal.Decimal `gorm:"type:numeric(20,8);default:0" json:"total_stock"`

	Variants    []Variant    `json:"variants"`
	SubVariants []SubVariant `json:"sub_variants"`
}
