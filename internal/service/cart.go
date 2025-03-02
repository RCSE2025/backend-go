package service

import (
	"context"
	"errors"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/repo"
)

type CartService struct {
	repo *repo.CartRepo
	pr   *repo.ProductRepo
}

func NewCartService(repo *repo.CartRepo, pr *repo.ProductRepo) *CartService {
	return &CartService{
		repo: repo,
		pr:   pr,
	}
}

func (s *CartService) PostInCart(userID, productID int64, quantity int) (model.CartItem, error) {
	_, err := s.pr.GetProductByID(context.Background(), productID)
	if err != nil {
		return model.CartItem{}, err
	}

	userCart, err := s.repo.GetCart(userID)
	if err != nil {
		return model.CartItem{}, err
	}
	for _, p := range userCart {
		if p.ProductID == productID {
			return model.CartItem{}, errors.New("product already in cart")
		}
	}
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

func (s *CartService) SetCartQuantity(userID, productID int64, quantity int) error {
	_, err := s.pr.GetProductByID(context.Background(), productID)
	if err != nil {
		return err
	}

	return s.repo.SetCartQuantity(userID, productID, quantity)
}
