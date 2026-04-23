package services

import (
	"github.com/m0xyu/learning-go-shop/internal/dto"
	"github.com/m0xyu/learning-go-shop/internal/models"
	"github.com/m0xyu/learning-go-shop/internal/repositories"
	"github.com/m0xyu/learning-go-shop/internal/utils"
)

var _ ProductServiceInterface = (*ProductService)(nil)

type ProductService struct {
	productRepo      repositories.ProductRepositoryInterface
	categoryRepo     repositories.CategoryRepositoryInterface
	productImageRepo repositories.ProductImageRepositoryInterface
}

func NewProductService(
	productRepo repositories.ProductRepositoryInterface,
	categoryRepo repositories.CategoryRepositoryInterface,
	productImageRepo repositories.ProductImageRepositoryInterface,
) *ProductService {
	return &ProductService{productRepo: productRepo, categoryRepo: categoryRepo, productImageRepo: productImageRepo}
}

func (s *ProductService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(&category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil

}

func (s *ProductService) GetCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i := range categories {
		response[i] = dto.CategoryResponse{
			ID:          categories[i].ID,
			Name:        categories[i].Name,
			Description: categories[i].Description,
			IsActive:    categories[i].IsActive,
			CreatedAt:   categories[i].CreatedAt,
			UpdatedAt:   categories[i].UpdatedAt,
		}
	}

	return response, nil
}

func (s *ProductService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {

	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (s *ProductService) DeleteCategory(id uint) error {
	return s.categoryRepo.Delete(id)
}

func (s *ProductService) GetProduct(id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := s.convertToProductResponse(product)
	return &response, nil
}

func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}

	if err := s.productRepo.Create(&product); err != nil {
		return nil, err
	}

	return s.GetProduct(product.ID)
}

func (s *ProductService) GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	products, total, err := s.productRepo.GetAll(offset, limit)
	if err != nil {
		return nil, nil, err
	}

	response := make([]dto.ProductResponse, len(products))
	for i := range products {
		response[i] = s.convertToProductResponse(&products[i])
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return response, meta, nil
}

func (s *ProductService) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return s.GetProduct(product.ID)
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}

func (s *ProductService) AddProductImage(productID uint, imageURL, altText string) error {
	count, err := s.productImageRepo.GetImagesCountByProductID(productID)
	if err != nil {
		return err
	}

	image := models.ProductImage{
		ProductID: productID,
		URL:       imageURL,
		AltText:   altText,
		IsPrimary: count == 0, // 最初の画像はプライマリにする
	}

	return s.productImageRepo.Create(&image)
}

func (s *ProductService) convertToProductResponse(product *models.Product) dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
			CreatedAt: product.Images[i].CreatedAt,
		}
	}

	return dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category: dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			IsActive:    product.Category.IsActive,
			CreatedAt:   product.Category.CreatedAt,
			UpdatedAt:   product.Category.UpdatedAt,
		},
		Images:    images,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
