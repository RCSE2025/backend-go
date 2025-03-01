package repo

import (
	"context"
	"github.com/RCSE2025/backend-go/internal/model"
	"gorm.io/gorm"
)

// ProductRepo представляет репозиторий для работы с продуктами
type ProductRepo struct {
	db *gorm.DB
}

// NewProductRepo создает новый экземпляр ProductRepo
func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

// GetProductByID возвращает продукт по его ID
func (r *ProductRepo) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	// Загружаем изображения
	var images []model.ProductImage
	if err := r.db.Where("product_id = ?", id).Find(&images).Error; err != nil {
		return nil, err
	}
	product.Images = images

	// Загружаем характеристики
	var specifications []model.ProductSpecification
	if err := r.db.Where("product_id = ?", id).Find(&specifications).Error; err != nil {
		return nil, err
	}
	product.Specifications = specifications

	return &product, nil
}

// GetProductReviews возвращает отзывы на продукт
func (r *ProductRepo) GetProductReviews(ctx context.Context, productID int64) ([]model.ProductReview, error) {
	var reviews []model.ProductReview
	if err := r.db.Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		return nil, err
	}

	for i := range reviews {
		var images = make([]model.ReviewImages, 0)
		if err := r.db.Where("review_id = ?", reviews[i].ID).Find(&images).Error; err != nil {
			return nil, err
		}
		for _, image := range images {
			reviews[i].Images = append(reviews[i].Images, image.URL)
		}
	}

	return reviews, nil
}

// AddProductReview добавляет новый отзыв на продукт
func (r *ProductRepo) AddProductReview(ctx context.Context, review model.ProductReview) (*model.ProductReview, error) {
	if err := r.db.Create(&review).Error; err != nil {
		return nil, err
	}

	// Получаем средний рейтинг
	var avgRating float64
	if err := r.db.Model(&model.ProductReview{}).
		Select("AVG(rating)").
		Where("product_id = ?", review.ProductID).
		Scan(&avgRating).Error; err != nil {
		return nil, err
	}

	// Получаем количество отзывов
	var reviewCount int64
	if err := r.db.Model(&model.ProductReview{}).
		Where("product_id = ?", review.ProductID).
		Count(&reviewCount).Error; err != nil {
		return nil, err
	}

	// Обновляем продукт
	if err := r.db.Model(&model.Product{}).
		Where("id = ?", review.ProductID).
		Updates(map[string]interface{}{
			"rating":       avgRating,
			"review_count": reviewCount,
		}).Error; err != nil {
		return nil, err
	}

	return &review, nil
}

// FilterProducts фильтрует продукты по заданным критериям
func (r *ProductRepo) FilterProducts(ctx context.Context, filters model.ProductFilters) ([]model.Product, error) {
	query := r.db.Model(&model.Product{})

	// Применяем фильтры
	if filters.SearchQuery != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filters.SearchQuery+"%", "%"+filters.SearchQuery+"%")
	}

	if len(filters.Categories) > 0 {
		query = query.Where("category IN ?", filters.Categories)
	}

	if filters.PriceRange != nil {
		query = query.Where("price BETWEEN ? AND ?", filters.PriceRange.Min, filters.PriceRange.Max)
	}

	if len(filters.Brands) > 0 {
		query = query.Where("brand IN ?", filters.Brands)
	}

	if filters.Rating > 0 {
		query = query.Where("rating >= ?", filters.Rating)
	}

	if filters.InStock != nil && *filters.InStock {
		query = query.Where("quantity > 0")
	}

	if filters.OnSale != nil && *filters.OnSale {
		query = query.Where("discount > 0")
	}

	// Сортировка
	switch filters.SortBy {
	case "price-asc":
		query = query.Order("price ASC")
	case "price-desc":
		query = query.Order("price DESC")
	case "rating":
		query = query.Order("rating DESC")
	case "newest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("id DESC")
	}

	var products []model.Product
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	// Загружаем изображения для каждого продукта
	for i := range products {
		var images []model.ProductImage
		if err := r.db.Where("product_id = ?", products[i].ID).Find(&images).Error; err != nil {
			return nil, err
		}
		products[i].Images = images
	}

	return products, nil
}

// CreateProduct создает новый продукт
func (r *ProductRepo) CreateProduct(ctx context.Context, product model.Product) (*model.Product, error) {
	// Начинаем транзакцию
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Создаем продукт
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Сохраняем изображения, если они есть
	if len(product.Images) > 0 {
		for i := range product.Images {
			product.Images[i].ProductID = product.ID
			if err := tx.Create(&product.Images[i]).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Сохраняем характеристики, если они есть
	if len(product.Specifications) > 0 {
		for i := range product.Specifications {
			product.Specifications[i].ProductID = product.ID
			if err := tx.Create(&product.Specifications[i]).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Завершаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// UpdateProduct обновляет существующий продукт
func (r *ProductRepo) UpdateProduct(ctx context.Context, product model.Product) (*model.Product, error) {
	// Начинаем транзакцию
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Обновляем продукт
	if err := tx.Model(&model.Product{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
		"business_id":        product.BusinessID,
		"price":              product.Price,
		"title":              product.Title,
		"description":        product.Description,
		"quantity":           product.Quantity,
		"discount":           product.Discount,
		"category":           product.Category,
		"brand":              product.Brand,
		"sku":                product.SKU,
		"estimated_delivery": product.EstimatedDelivery,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Обновляем изображения, если они есть
	if len(product.Images) > 0 {
		// Удаляем старые изображения
		if err := tx.Where("product_id = ?", product.ID).Delete(&model.ProductImage{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Добавляем новые изображения
		for i := range product.Images {
			product.Images[i].ProductID = product.ID
			if err := tx.Create(&product.Images[i]).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Обновляем характеристики, если они есть
	if len(product.Specifications) > 0 {
		// Удаляем старые характеристики
		if err := tx.Where("product_id = ?", product.ID).Delete(&model.ProductSpecification{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Добавляем новые характеристики
		for i := range product.Specifications {
			product.Specifications[i].ProductID = product.ID
			if err := tx.Create(&product.Specifications[i]).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Завершаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// DeleteProduct удаляет продукт по его ID
func (r *ProductRepo) DeleteProduct(ctx context.Context, id int64) error {
	// Начинаем транзакцию
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Удаляем изображения
	if err := tx.Where("product_id = ?", id).Delete(&model.ProductImage{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Удаляем характеристики
	if err := tx.Where("product_id = ?", id).Delete(&model.ProductSpecification{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Удаляем отзывы
	if err := tx.Where("product_id = ?", id).Delete(&model.ProductReview{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Удаляем продукт
	if err := tx.Delete(&model.Product{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Завершаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// AddProductImage добавляет изображение продукта
func (r *ProductRepo) AddProductImage(ctx context.Context, image model.ProductImage) (*model.ProductImage, error) {
	// Если это основное изображение, сбрасываем флаг у других изображений
	if image.IsPrimary {
		if err := r.db.Model(&model.ProductImage{}).
			Where("product_id = ?", image.ProductID).
			Update("is_primary", false).Error; err != nil {
			return nil, err
		}
	}

	// Создаем изображение
	if err := r.db.Create(&image).Error; err != nil {
		return nil, err
	}

	return &image, nil
}

// DeleteProductImage удаляет изображение продукта
func (r *ProductRepo) DeleteProductImage(ctx context.Context, imageID int64) error {
	// Получаем информацию об изображении перед удалением
	var image model.ProductImage
	if err := r.db.First(&image, imageID).Error; err != nil {
		return err
	}

	// Удаляем изображение
	if err := r.db.Delete(&model.ProductImage{}, imageID).Error; err != nil {
		return err
	}

	// Если удаленное изображение было основным, назначаем новое основное изображение
	if image.IsPrimary {
		var newPrimaryImage model.ProductImage
		if err := r.db.Where("product_id = ?", image.ProductID).First(&newPrimaryImage).Error; err == nil {
			// Если нашли другое изображение, делаем его основным
			if err := r.db.Model(&newPrimaryImage).Update("is_primary", true).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// GetProductImages возвращает изображения продукта
func (r *ProductRepo) GetProductImages(ctx context.Context, productID int64) ([]model.ProductImage, error) {
	var images []model.ProductImage
	if err := r.db.Where("product_id = ?", productID).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *ProductRepo) GetReviewImages(reviewID int64) ([]model.ReviewImages, error) {
	var images []model.ReviewImages
	if err := r.db.Where("review_id = ?", reviewID).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *ProductRepo) UploadReviewImages(image model.ReviewImages) (model.ReviewImages, error) {
	if err := r.db.Create(&image).Error; err != nil {
		return model.ReviewImages{}, err
	}
	return image, nil
}
