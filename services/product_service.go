package services

import (
	"errors"

	"inventory-backend/config"
	"inventory-backend/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ---------- REQUEST DTOs ----------

type VariantRequest struct {
	Name    string   `json:"name" binding:"required"`
	Options []string `json:"options" binding:"required"`
}

type SubVariantOptionValue struct {
	VariantName string `json:"variant_name" binding:"required"`
	Value       string `json:"value" binding:"required"`
}

type SubVariantRequest struct {
	SKU          string                  `json:"sku" binding:"required"`
	OptionValues []SubVariantOptionValue `json:"option_values" binding:"required"`
}

type CreateProductRequest struct {
	ProductID    int64               `json:"product_id" binding:"required"`
	ProductCode  string              `json:"product_code" binding:"required"`
	ProductName  string              `json:"product_name" binding:"required"`
	ProductImage string              `json:"product_image"`
	CreatedUser  string              `json:"created_user" binding:"required"`
	IsFavourite  bool                `json:"is_favourite"`
	Active       bool                `json:"active"`
	HSNCode      string              `json:"hsn_code"`
	Variants     []VariantRequest    `json:"variants"`
	SubVariants  []SubVariantRequest `json:"sub_variants"`
}

// ---------- SERVICE IMPLEMENTATIONS ----------

func CreateProduct(req CreateProductRequest) (*models.Product, error) {
	db := config.DB

	createdUserUUID, err := uuid.Parse(req.CreatedUser)
	if err != nil {
		return nil, errors.New("invalid created_user UUID")
	}

	var product models.Product

	err = db.Transaction(func(tx *gorm.DB) error {
		product = models.Product{
			ID:           uuid.New(),
			ProductID:    req.ProductID,
			ProductCode:  req.ProductCode,
			ProductName:  req.ProductName,
			ProductImage: req.ProductImage,
			CreatedUser:  createdUserUUID,
			IsFavourite:  req.IsFavourite,
			Active:       req.Active,
			HSNCode:      req.HSNCode,
			TotalStock:   decimal.NewFromInt(0),
		}

		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		// Maps to resolve variant/option IDs by names
		variantMap := make(map[string]uuid.UUID)
		optionMap := make(map[string]uuid.UUID) // key: variantName|value

		// Create variants + options
		for _, v := range req.Variants {
			variant := models.Variant{
				ID:        uuid.New(),
				ProductID: product.ID,
				Name:      v.Name,
			}
			if err := tx.Create(&variant).Error; err != nil {
				return err
			}
			variantMap[v.Name] = variant.ID

			for _, optVal := range v.Options {
				option := models.VariantOption{
					ID:        uuid.New(),
					VariantID: variant.ID,
					Value:     optVal,
				}
				if err := tx.Create(&option).Error; err != nil {
					return err
				}
				key := v.Name + "|" + optVal
				optionMap[key] = option.ID
			}
		}

		// Create sub-variants
		for _, sv := range req.SubVariants {
			var optionIDs []string
			for _, ov := range sv.OptionValues {
				key := ov.VariantName + "|" + ov.Value
				optID, ok := optionMap[key]
				if !ok {
					return errors.New("sub_variant references unknown variant/option: " + key)
				}
				optionIDs = append(optionIDs, optID.String())
			}

			subVariant := models.SubVariant{
				ID:        uuid.New(),
				ProductID: product.ID,
				OptionIDs: optionIDs,
				SKU:       sv.SKU,
				Stock:     decimal.NewFromInt(0),
			}
			if err := tx.Create(&subVariant).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload with relations
	if err := db.Preload("Variants.Options").
		Preload("SubVariants").
		First(&product, "id = ?", product.ID).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func ListProducts(page, limit int) ([]models.Product, int64, error) {
	db := config.DB

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var products []models.Product
	var total int64

	if err := db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Variants.Options").
		Preload("SubVariants").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
