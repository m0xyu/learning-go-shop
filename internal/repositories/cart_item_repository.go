package repositories

import (
	"errors"

	"github.com/m0xyu/learning-go-shop/internal/models"
	"gorm.io/gorm"
)

type CartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) *CartItemRepository {
	return &CartItemRepository{db: db}
}

func (r *CartItemRepository) GetByCartIDAndProductID(cartID uint, productID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	if err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).Find(&cartItems).Error; err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (r *CartItemRepository) GetByUserIDAndItemID(userID, itemID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	if err := r.db.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&cartItem).Error; err != nil {
		return nil, errors.New("cart item not found")
	}
	return &cartItem, nil
}

func (r *CartItemRepository) GetByCartAndProduct(cartID, productID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.Unscoped().Where("cart_id = ? AND product_id = ?", cartID, productID).First(&cartItem).Error
	if err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *CartItemRepository) Create(cartItem *models.CartItem) error {
	return r.db.Create(cartItem).Error
}

func (r *CartItemRepository) Update(cartItem *models.CartItem) error {
	return r.db.Save(cartItem).Error
}

func (r *CartItemRepository) Delete(userID, itemID uint) error {
	return r.db.Where("id = ? AND cart_id IN (?)", itemID,
		r.db.Select("id").Table("carts").
			Where("user_id = ?", userID)).
		Delete(&models.CartItem{}).Error
}
