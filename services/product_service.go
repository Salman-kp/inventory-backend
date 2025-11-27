package services

import (
	"inventory-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) CreateProduct(p *models.Product) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return s.db.Create(p).Error
}

func (s *ProductService) ListProducts(page, limit int) ([]models.Product, error) {
	var products []models.Product

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	err := s.db.
		Limit(limit).
		Offset(offset).
		Order("created_date DESC").
		Find(&products).Error

	return products, err
}
