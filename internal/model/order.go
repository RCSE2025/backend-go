package model

type Order struct {
	BaseModel
	ID             int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int64           `json:"user_id" gorm:"not null"`
	Status         OrderStatusType `json:"status" gorm:"not null;default:created" swaggertype:"primitive,string"`
	PaymentConfirm bool            `json:"payment_confirm" gorm:"not null;default:false"`
	Address        string          `json:"address" gorm:"not null"`
}

type OrderStatusType string

const StatusCreated OrderStatusType = "created"
const StatusDelivery OrderStatusType = "delivery"
const StatusClosed OrderStatusType = "closed"

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	BaseModel
	UserID    int64   `json:"user_id" gorm:"not null"`
	OrderID   int64   `json:"order_id" gorm:"not null"`
	ProductID int64   `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
}

type ExtendedOrderItem struct {
	OrderItem `json:"order_item"`
	Product   `json:"product"`
}

type OrderItemResponse struct {
	Order
	OrderItems []ExtendedOrderItem `json:"order_items"`
}
