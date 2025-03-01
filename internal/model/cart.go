package model

type CartItem struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (CartItem) TableName() string { return "cart_items" }

type CartItemsResponse struct {
	CartItem
	Product Product `json:"product"`
}
