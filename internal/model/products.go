package model

//
//- id 		PK[int]
//- business_id	FK[int] -> Businesses
//- price		Decimal or float
//- title		str
//- description 	str
//- quantity	int

type Product struct {
	baseModel
	ID          int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	BusinessID  int64   `json:"business_id" gorm:"not null"`
	Price       float64 `json:"price" gorm:"not null"`
	Title       string  `json:"title" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Quantity    int     `json:"quantity" gorm:"not null"`
}

func (Product) TableName() string {
	return "products"
}

type ProductImage struct {
	baseModel
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID int64  `json:"product_id" gorm:"not null"`
	FileUUID  string `json:"file_uuid" gorm:"not null"`
}

func (ProductImage) TableName() string {
	return "product_images"
}

type ProductWithImages struct {
	Product
	Images []string `json:"images"`
}
