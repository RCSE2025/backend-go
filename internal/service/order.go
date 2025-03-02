package service

import (
	"context"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/repo"
)

type OrderService struct {
	repo        *repo.OrderRepo
	productRepo *repo.ProductRepo
	yookassa    *YookassaPayment
}

func NewOrderService(repo *repo.OrderRepo, productRepo *repo.ProductRepo, yookassa *YookassaPayment) *OrderService {
	return &OrderService{
		repo:        repo,
		productRepo: productRepo,
		yookassa:    yookassa,
	}
}

func (ordS *OrderService) CreateOrder(userID int64) (model.Order, error) {
	return ordS.repo.CreateOrder(model.Order{UserID: userID, Status: model.StatusCreated})
}

func (ordS *OrderService) CreateOrderItem(userID, orderID, productID int64, quantity int) (model.OrderItem, error) {
	product, err := ordS.productRepo.GetProductByID(context.Background(), productID)
	if err != nil {
		return model.OrderItem{}, err
	}

	return ordS.repo.CreateOrderItem(model.OrderItem{UserID: userID, OrderID: orderID, ProductID: productID, Quantity: quantity, Price: product.Price})
}

func (ordS *OrderService) SetOrderStatus(orderID, userID int64, status model.OrderStatusType) error {
	return ordS.repo.SetOrderStatus(orderID, userID, status)
}

func (ordS *OrderService) GetUserOrders(userID int64) ([]model.OrderItemResponse, error) {
	return ordS.repo.GetUserOrders(userID)
}

func (ordS *OrderService) ConfirmOrderPayment(orderID int64) error {
	return ordS.repo.ConfirmOrderPayment(orderID)
}

func (ordS *OrderService) CreateOrderPayment(orderID int64, amount float64) (string, error) {
	return ordS.yookassa.CreateOrderPayment(orderID, amount)
}
