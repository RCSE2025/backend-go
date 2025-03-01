package repo

import (
	"github.com/RCSE2025/backend-go/internal/model"
	"gorm.io/gorm"
)

type CartRepo struct {
	db *gorm.DB
}

func NewCartRepo(db *gorm.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) PostInCart(item model.CartItem) (model.CartItem, error) {
	return item, r.db.Create(&item).Error
}

func (r *CartRepo) DeleteFromCart(userID int64, ids []int64) error {
	return r.db.Where("id IN ? AND user_id = ?", ids, userID).Delete(&model.CartItem{}).Error
}

func (r *CartRepo) GetCart(userID int64) ([]model.CartItemsResponse, error) {
	var carts []model.CartItemsResponse
	return carts, r.db.Joins("JOIN products ON products.id = product_id").Where("user_id = ?", userID).Find(&carts).Error
}
