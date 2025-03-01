package repo

import (
	"context"
	"github.com/RCSE2025/backend-go/internal/model"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db       *gorm.DB
	prodRepo *ProductRepo
}

func NewOrderRepo(db *gorm.DB, prodRepo *ProductRepo) *OrderRepo {
	return &OrderRepo{db: db, prodRepo: prodRepo}
}

func (or *OrderRepo) CreateOrder(order model.Order) (model.Order, error) {
	return order, or.db.Create(&order).Error
}

func (or *OrderRepo) CreateOrderItem(orderItem model.OrderItem) (model.OrderItem, error) {
	return orderItem, or.db.Create(&orderItem).Error
}

func (or *OrderRepo) SetOrderStatus(userID, orderID int64, status string) error {
	return or.db.Model(&model.Order{}).Where("user_id = ? AND id = ?", userID, orderID).Update("status", status).Error
}

func (or *OrderRepo) GetUserOrders(userID int64) ([]model.OrderItemResponse, error) {
	var userOrders = make([]model.OrderItemResponse, 0)

	var orders = make([]model.Order, 0)
	err := or.db.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return []model.OrderItemResponse{}, err
	}

	for _, order := range orders {
		userOrder := model.OrderItemResponse{Order: order, OrderItems: []model.ExtendedOrderItem{}}
		var orderItems = make([]model.OrderItem, 0)
		err := or.db.Where("order_id = ? AND user_id = ?", order.ID, userID).Find(&orderItems).Error
		if err != nil {
			return []model.OrderItemResponse{}, err
		}

		for _, orderItem := range orderItems {
			product, err := or.prodRepo.GetProductByID(context.Background(), orderItem.ProductID)
			if err != nil {
				return []model.OrderItemResponse{}, err
			}

			userOrder.OrderItems = append(userOrder.OrderItems, model.ExtendedOrderItem{OrderItem: orderItem, Product: *product})
		}
		userOrders = append(userOrders, userOrder)
	}
	return userOrders, nil
}
