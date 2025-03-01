package repo

import (
	"context"
	"github.com/RCSE2025/backend-go/internal/model"
	"gorm.io/gorm"
)

type CartRepo struct {
	db *gorm.DB
	pr *ProductRepo
}

func NewCartRepo(db *gorm.DB, pr *ProductRepo) *CartRepo {
	return &CartRepo{db: db, pr: pr}
}

func (r *CartRepo) PostInCart(item model.CartItem) (model.CartItem, error) {
	return item, r.db.Create(&item).Error
}

func (r *CartRepo) DeleteFromCart(userID int64, ids []int64) error {
	return r.db.Where("product_id IN ? AND user_id = ?", ids, userID).Delete(&model.CartItem{}).Error
}

func (r *CartRepo) GetCart(userID int64) ([]model.CartItemsResponse, error) {
	var cartItems []model.CartItem
	if err := r.db.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	var result []model.CartItemsResponse
	for _, item := range cartItems {
		product, err := r.pr.GetProductByID(context.Background(), item.ProductID)
		if err != nil {
			return nil, err
		}

		result = append(result, model.CartItemsResponse{
			CartItem: item,
			Product:  *product,
		})
	}

	return result, nil
}

func (r *CartRepo) SetCartQuantity(userID, productID int64, quantity int) error {
	return r.db.Model(&model.CartItem{}).Where("user_id = ? AND product_id = ?", userID, productID).Update("quantity", quantity).Error
}
