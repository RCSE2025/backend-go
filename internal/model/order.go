package model

//Orders >>
//- id			PK[int]
//- user_id 	FK[int]
//- status		Enum(created, delivery, closed)
//
//
//OrderItems >>
//- user_id		PK,FK[int]
//- product_id	FK[int]
//- quantity	int
//- price 		Decimal or float

type Order struct {
	baseModel
	ID     int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID int64 `json:"user_id" gorm:"not null"`
	Status string
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	baseModel
	OrderID   int64   `json:"order_id" gorm:"not null"`
	ProductID int64   `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
}
