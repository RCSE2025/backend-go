package model

import "time"

// ProductCategory представляет категорию товара
type ProductCategory string

// Константы для категорий товаров
const (
	ProductCategoryElectronics ProductCategory = "ELECTRONICS"
	ProductCategoryHome        ProductCategory = "HOME"
	ProductCategoryFashion     ProductCategory = "FASHION"
	ProductCategorySports      ProductCategory = "SPORTS"
	ProductCategoryBeauty      ProductCategory = "BEAUTY"
	ProductCategoryToys        ProductCategory = "TOYS"
	ProductCategoryBooks       ProductCategory = "BOOKS"
	ProductCategoryFood        ProductCategory = "FOOD"
	ProductCategoryOther       ProductCategory = "OTHER"
)

// ProductCategoryMap маппинг категорий на названия
var ProductCategoryMap = map[ProductCategory]string{
	ProductCategoryElectronics: "Электроника",
	ProductCategoryHome:        "Дом и кухня",
	ProductCategoryFashion:     "Мода",
	ProductCategorySports:      "Спорт и отдых",
	ProductCategoryBeauty:      "Красота и здоровье",
	ProductCategoryToys:        "Игрушки",
	ProductCategoryBooks:       "Книги",
	ProductCategoryFood:        "Продукты питания",
	ProductCategoryOther:       "Другое",
}

// ProductSpecification представляет характеристику товара
type ProductSpecification struct {
	BaseModel
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID int64  `json:"product_id" gorm:"not null"`
	Name      string `json:"name" gorm:"not null"`
	Value     string `json:"value" gorm:"not null"`
}

func (ProductSpecification) TableName() string {
	return "product_specifications"
}

// ProductReview представляет отзыв на товар
type ProductReview struct {
	BaseModel
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID int64     `json:"product_id" gorm:"not null"`
	UserID    int64     `json:"user_id" gorm:"not null"`
	UserName  string    `json:"user_name" gorm:"not null"`
	Rating    int       `json:"rating" gorm:"not null"`
	Comment   string    `json:"comment" gorm:"not null"`
	Date      time.Time `json:"date" gorm:"not null"`
	Images    []string  `json:"images" gorm:"-"` // Не хранится в базе данных напрямую
}

func (ProductReview) TableName() string {
	return "product_reviews"
}

type ReviewImages struct {
	BaseModel
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ReviewID  int64  `json:"review_id" gorm:"not null"`
	FileUUID  string `json:"file_uuid" gorm:"not null"`
	URL       string `json:"url" gorm:"null"`
	IsPrimary bool   `json:"is_primary" gorm:"default:false"`
}

func (ReviewImages) TableName() string { return "review_images" }

// ProductImage представляет изображение товара
type ProductImage struct {
	BaseModel
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID int64  `json:"product_id" gorm:"not null"`
	FileUUID  string `json:"file_uuid" gorm:"not null"`
	URL       string `json:"url" gorm:"null"` // Не хранится в базе данных напрямую
	IsPrimary bool   `json:"is_primary" gorm:"default:false"`
}

func (ProductImage) TableName() string {
	return "product_images"
}

// SetTimestamps устанавливает время создания и обновления
func (p *ProductImage) SetTimestamps() {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
}

func (r *ReviewImages) SetTimestamps() {
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
}

type ProductStatus string

const StatusConsideration ProductStatus = "consideration"
const StatusReject ProductStatus = "reject"
const StatusApprove ProductStatus = "approve"

// Product представляет товар
type Product struct {
	BaseModel
	ID                int64                  `json:"id" gorm:"primaryKey;autoIncrement"`
	BusinessID        int64                  `json:"business_id" gorm:"not null"`
	Price             float64                `json:"price" gorm:"not null"`
	Title             string                 `json:"title" gorm:"not null"`
	Description       string                 `json:"description" gorm:"not null"`
	Quantity          int                    `json:"quantity" gorm:"not null"`
	Rating            float64                `json:"rating" gorm:"default:0"`
	ReviewCount       int                    `json:"review_count" gorm:"default:0"`
	Discount          float64                `json:"discount" gorm:"default:0"`
	Category          ProductCategory        `json:"category" gorm:"type:varchar(50);default:'OTHER'"`
	Brand             string                 `json:"brand" gorm:"default:''"`
	SKU               string                 `json:"sku" gorm:"default:''"`
	EstimatedDelivery string                 `json:"estimated_delivery" gorm:"default:'3-5 дней'"`
	Images            []ProductImage         `json:"images" gorm:"-"`           // Загружается отдельно
	Specifications    []ProductSpecification `json:"specifications" gorm:"-"`   // Загружается отдельно
	Reviews           []ProductReview        `json:"reviews" gorm:"-"`          // Загружается отдельно
	RelatedProducts   []int64                `json:"related_products" gorm:"-"` // Загружается отдельно

	ProductStatus `json:"status" gorm:"not null;default:'consideration'"`
}

func (Product) TableName() string {
	return "products"
}

// PriceRange представляет диапазон цен для фильтрации
type PriceRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// CategoryFilter представляет фильтр по категории
type CategoryFilter struct {
	ID       string          `json:"id"`
	Title    string          `json:"title"`
	Image    string          `json:"image"`
	Link     string          `json:"link"`
	Category ProductCategory `json:"category"`
}

// ProductFilters представляет фильтры для поиска товаров
type ProductFilters struct {
	SearchQuery string            `json:"search_query,omitempty"`
	Categories  []ProductCategory `json:"categories,omitempty"`
	PriceRange  *PriceRange       `json:"price_range,omitempty"`
	Brands      []string          `json:"brands,omitempty"`
	Rating      float64           `json:"rating,omitempty"`
	InStock     *bool             `json:"in_stock,omitempty"`
	OnSale      *bool             `json:"on_sale,omitempty"`
	SortBy      string            `json:"sort_by,omitempty"` // price-asc, price-desc, rating, newest
}

type ProductQueryParams struct {
	SearchQuery string            `form:"q"`
	Categories  []ProductCategory `form:"categories"`
	MinPrice    float64           `form:"min_price"`
	MaxPrice    float64           `form:"max_price"`
	Brands      []string          `form:"brands"`
	Rating      float64           `form:"rating"`
	InStock     *bool             `form:"in_stock"`
	OnSale      *bool             `form:"on_sale"`
	SortBy      string            `form:"sort_by"` // price-asc, price-desc, rating, newest
	Page        int               `form:"page,default=1"`
	PageSize    int               `form:"per_page,default=20"`
}

// ProductCreateRequest представляет данные для создания нового продукта
type ProductCreateRequest struct {
	BusinessID        int64                  `json:"business_id" binding:"required"`
	Price             float64                `json:"price" binding:"required,gt=0"`
	Title             string                 `json:"title" binding:"required"`
	Description       string                 `json:"description" binding:"required"`
	Quantity          int                    `json:"quantity" binding:"required,gte=0"`
	Discount          float64                `json:"discount" binding:"omitempty,gte=0"`
	Category          ProductCategory        `json:"category" binding:"required"`
	Brand             string                 `json:"brand" binding:"omitempty"`
	SKU               string                 `json:"sku" binding:"omitempty"`
	EstimatedDelivery string                 `json:"estimated_delivery" binding:"omitempty"`
	Specifications    []ProductSpecification `json:"specifications" binding:"omitempty"`
}

// ToProduct преобразует ProductCreateRequest в Product
func (r *ProductCreateRequest) ToProduct() Product {
	return Product{
		BusinessID:        r.BusinessID,
		Price:             r.Price,
		Title:             r.Title,
		Description:       r.Description,
		Quantity:          r.Quantity,
		Discount:          r.Discount,
		Category:          r.Category,
		Brand:             r.Brand,
		SKU:               r.SKU,
		EstimatedDelivery: r.EstimatedDelivery,
		Specifications:    r.Specifications,
		Rating:            0,
		ReviewCount:       0,
	}
}

// ProductUpdateRequest представляет данные для обновления продукта
type ProductUpdateRequest struct {
	BusinessID        int64                  `json:"business_id" binding:"omitempty"`
	Price             float64                `json:"price" binding:"omitempty,gt=0"`
	Title             string                 `json:"title" binding:"omitempty"`
	Description       string                 `json:"description" binding:"omitempty"`
	Quantity          int                    `json:"quantity" binding:"omitempty,gte=0"`
	Discount          float64                `json:"discount" binding:"omitempty,gte=0"`
	Category          ProductCategory        `json:"category" binding:"omitempty"`
	Brand             string                 `json:"brand" binding:"omitempty"`
	SKU               string                 `json:"sku" binding:"omitempty"`
	EstimatedDelivery string                 `json:"estimated_delivery" binding:"omitempty"`
	Specifications    []ProductSpecification `json:"specifications" binding:"omitempty"`
}

// ApplyToProduct применяет изменения из ProductUpdateRequest к Product
func (r *ProductUpdateRequest) ApplyToProduct(product *Product) {
	if r.BusinessID != 0 {
		product.BusinessID = r.BusinessID
	}
	if r.Price > 0 {
		product.Price = r.Price
	}
	if r.Title != "" {
		product.Title = r.Title
	}
	if r.Description != "" {
		product.Description = r.Description
	}
	if r.Quantity >= 0 {
		product.Quantity = r.Quantity
	}
	if r.Discount >= 0 {
		product.Discount = r.Discount
	}
	if r.Category != "" {
		product.Category = r.Category
	}
	if r.Brand != "" {
		product.Brand = r.Brand
	}
	if r.SKU != "" {
		product.SKU = r.SKU
	}
	if r.EstimatedDelivery != "" {
		product.EstimatedDelivery = r.EstimatedDelivery
	}
	if len(r.Specifications) > 0 {
		product.Specifications = r.Specifications
	}
}
