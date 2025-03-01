package service

import (
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/repo"
)

type CartService struct {
	repo *repo.CartRepo
}

func NewCartService(repo *repo.CartRepo) *CartService {
	return &CartService{
		repo: repo,
	}
}

func (s *CartService) PostInCart(userID, productID int64, quantity int) (model.CartItem, error) {
	cart, err := s.repo.PostInCart(model.CartItem{UserID: userID, Quantity: quantity, ProductID: productID})
	return cart, err
}

func (s *CartService) DeleteCart(userID int64, productIDs []int64) error {
	return s.repo.DeleteFromCart(userID, productIDs)
}

func (s *CartService) GetUserCart(userID int64) ([]model.CartItemsResponse, error) {
	cards, err := s.repo.GetCart(userID)
	return cards, err
}
