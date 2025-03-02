package service

import (
	"context"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/repo"
	"github.com/RCSE2025/backend-go/internal/utils"
)

// ProductService представляет сервис для работы с продуктами
type ProductService struct {
	repo           *repo.ProductRepo
	s3Worker       *utils.S3WorkerAPI
	s3WorkerReview *utils.S3WorkerAPI
}

// NewProductService создает новый экземпляр ProductService
func NewProductService(repo *repo.ProductRepo, s3Worker, s3WorkerReview *utils.S3WorkerAPI) *ProductService {
	return &ProductService{
		repo:           repo,
		s3Worker:       s3Worker,
		s3WorkerReview: s3WorkerReview,
	}
}

// GetS3Worker возвращает S3 Worker API
func (s *ProductService) GetS3Worker() *utils.S3WorkerAPI {
	return s.s3Worker
}

func (s *ProductService) GetS3WorkerReview() *utils.S3WorkerAPI {
	return s.s3WorkerReview
}

// GetProductByID возвращает продукт по его ID
func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

// GetProductReviews возвращает отзывы на продукт
func (s *ProductService) GetProductReviews(ctx context.Context, productID int64) ([]model.ProductReview, error) {
	return s.repo.GetProductReviews(ctx, productID)
}

// AddProductReview добавляет новый отзыв на продукт
func (s *ProductService) AddProductReview(ctx context.Context, review model.ProductReview) (*model.ProductReview, error) {
	return s.repo.AddProductReview(ctx, review)
}

// FilterProducts фильтрует продукты по заданным критериям
func (s *ProductService) FilterProducts(ctx context.Context, filters model.ProductQueryParams) ([]model.Product, error) {
	return s.repo.FilterProducts(ctx, filters)
}

// CreateProduct создает новый продукт
func (s *ProductService) CreateProduct(ctx context.Context, product model.Product) (*model.Product, error) {
	return s.repo.CreateProduct(ctx, product)
}

// UpdateProduct обновляет существующий продукт
func (s *ProductService) UpdateProduct(ctx context.Context, product model.Product) (*model.Product, error) {
	return s.repo.UpdateProduct(ctx, product)
}

// DeleteProduct удаляет продукт по его ID
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}

// AddProductImage добавляет изображение продукта
func (s *ProductService) AddProductImage(ctx context.Context, image model.ProductImage) (*model.ProductImage, error) {
	return s.repo.AddProductImage(ctx, image)
}

// DeleteProductImage удаляет изображение продукта
func (s *ProductService) DeleteProductImage(ctx context.Context, imageID int64) error {
	return s.repo.DeleteProductImage(ctx, imageID)
}

// GetProductImages возвращает изображения продукта
func (s *ProductService) GetProductImages(ctx context.Context, productID int64) ([]model.ProductImage, error) {
	return s.repo.GetProductImages(ctx, productID)
}

// GetProductCategories возвращает все категории продуктов
func (s *ProductService) GetProductCategories() []model.CategoryFilter {
	categories := make([]model.CategoryFilter, 0, len(model.ProductCategoryMap))

	for category, title := range model.ProductCategoryMap {
		categoryFilter := model.CategoryFilter{
			ID:       string(category),
			Title:    title,
			Image:    "/images/categories/" + string(category) + ".jpg", // Пример пути к изображению
			Link:     "/catalog?categories=" + string(category),         // Пример ссылки на категорию
			Category: category,
		}
		categories = append(categories, categoryFilter)
	}

	return categories
}

func (s *ProductService) AddReviewImage(image model.ReviewImages) (model.ReviewImages, error) {
	return s.repo.UploadReviewImages(image)
}

func (s *ProductService) GetReviewImages(reviewID int64) ([]model.ReviewImages, error) {
	return s.repo.GetReviewImages(reviewID)
}
