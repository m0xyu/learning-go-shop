package repositories

import (
	"github.com/m0xyu/learning-go-shop/internal/models"
	"gorm.io/gorm"
)

type ProductImageRepository struct {
	db *gorm.DB
}

func NewProductImageRepository(db *gorm.DB) *ProductImageRepository {
	return &ProductImageRepository{db: db}
}

func (r *ProductImageRepository) GetImagesCountByProductID(productID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ProductImageRepository) Create(image *models.ProductImage) error {
	return r.db.Create(image).Error
}
