package repositories

import (
	"github.com/m0xyu/learning-go-shop/internal/dto"
	"github.com/m0xyu/learning-go-shop/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAll(offset, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// 件数の取得
	if err := r.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// データ本体の取得
	if err := r.db.Preload("Category").Preload("Images").
		Where("is_active = ?", true).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) Search(req *dto.SearchProductsRequest) ([]models.ProductsWithRank, int64, error) {
	var total int64

	offset := (req.Page - 1) * req.Limit

	// クエリの組み立て（ここに提示されたロジックを書く）
	query := r.db.Model(&models.Product{}).
		Select("products.*, ts_rank(search_vector, plainto_tsquery('english', ?)) as rank", req.Query).
		Where("search_vector @@ plainto_tsquery('english', ?)", req.Query).
		Where("is_active = ?", true)

	// ... フィルタ条件の追加 ...
	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	if req.MinPrice != nil {
		query = query.Where("price >= ?", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		query = query.Where("price <= ?", *req.MaxPrice)
	}

	query.Count(&total)

	// Execute query with ranking and create product slices
	var rows []models.ProductsWithRank
	err := query.
		Order("rank DESC, created_at DESC"). // order by relevance
		Preload("Category").
		Preload("Images").
		Offset(offset).
		Limit(req.Limit).
		Find(&rows).Error

	return rows, total, err
}
