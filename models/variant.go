package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Variant struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID uuid.UUID       `gorm:"type:uuid;index" json:"product_id"`
	Name      string          `json:"name"`
	Options   []VariantOption `json:"options"`
}

type VariantOption struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	VariantID uuid.UUID `gorm:"type:uuid;index" json:"variant_id"`
	Value     string    `json:"value"`
}

type SubVariant struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID uuid.UUID       `gorm:"type:uuid;index" json:"product_id"`
	OptionIDs pq.StringArray  `gorm:"type:text[]" json:"option_ids"`
	SKU       string          `json:"sku"`
	Stock     decimal.Decimal `gorm:"type:numeric(20,8);default:0" json:"stock"`
}