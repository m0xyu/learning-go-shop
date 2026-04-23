package services

import (
	"errors"
	"fmt"

	"github.com/m0xyu/learning-go-shop/internal/dto"
	"github.com/m0xyu/learning-go-shop/internal/models"
	"github.com/m0xyu/learning-go-shop/internal/repositories"
	"gorm.io/gorm"
)

var _ CartServiceInterface = (*CartService)(nil)

type CartService struct {
	productRepo  repositories.ProductRepositoryInterface
	cartRepo     repositories.CartRepositoryInterface
	cartItemRepo repositories.CartItemRepositoryInterface
}

func NewCartService(
	productRepo repositories.ProductRepositoryInterface,
	cartRepo repositories.CartRepositoryInterface,
	cartItemRepo repositories.CartItemRepositoryInterface,
) *CartService {
	return &CartService{productRepo: productRepo, cartRepo: cartRepo, cartItemRepo: cartItemRepo}
}

func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return s.convertToCartResponse(cart), nil
}

func (s *CartService) AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		cart = &models.Cart{UserID: userID}
		if err := s.cartRepo.Create(cart); err != nil {
			return nil, err
		}
	}

	cartItem, err := s.cartItemRepo.GetByCartAndProduct(cart.ID, req.ProductID)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cartItem = &models.CartItem{
				CartID:    cart.ID,
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			}
			if err := s.cartItemRepo.Create(cartItem); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		fmt.Printf("CartItem: %+v, err: %v\n", cartItem, err)
		cartItem.Quantity += req.Quantity
		if err := s.cartItemRepo.Update(cartItem); err != nil {
			return nil, err
		}
	}

	return s.GetCart(userID)
}

func (s *CartService) UpdateCartItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	cartItem, err := s.cartItemRepo.GetByUserIDAndItemID(userID, itemID)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	product, err := s.productRepo.GetByID(cartItem.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	cartItem.Quantity = req.Quantity
	if err := s.cartItemRepo.Update(cartItem); err != nil {
		return nil, err
	}

	return s.GetCart(userID)
}

func (s *CartService) RemoveFromCart(userID, itemID uint) error {
	return s.cartItemRepo.Delete(userID, itemID)
}

func (s *CartService) convertToCartResponse(cart *models.Cart) *dto.CartResponse {

	cartItems := make([]dto.CartItemResponse, len(cart.CartItems)) // memory allocation
	var total float64

	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
			Quantity:  cart.CartItems[i].Quantity,
			Subtotal:  subtotal,
			CreatedAt: cart.CartItems[i].CreatedAt,
			UpdatedAt: cart.CartItems[i].UpdatedAt,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
