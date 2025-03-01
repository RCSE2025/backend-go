package product

import (
	"github.com/RCSE2025/backend-go/internal/http/middleware/logger"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type productRoutes struct {
	productService *service.ProductService
}

func NewProductRoutes(h *gin.RouterGroup, jwtService service.JWTService, productService *service.ProductService) {
	g := h.Group("/product")

	pr := productRoutes{
		productService: productService,
	}

	g.POST("/images/upload", pr.uploadImages)
	g.GET("/categories", pr.getCategories)
	g.GET("/:id", pr.getProduct)
	g.GET("/:id/reviews", pr.getProductReviews)
	g.POST("/:id/reviews", pr.addProductReview)
	g.GET("/filter", pr.filterProducts)
	g.POST("", pr.createProduct)
	g.PUT("/:id", pr.updateProduct)
	g.DELETE("/:id", pr.deleteProduct)
	g.POST("/:id/images/upload", pr.UploadReviewFile)
}

// uploadImages
// @Summary     Upload multiple images
// @Description Upload multiple images
// @Tags  	    product
// @Accept      multipart/form-data
// @Produce     json
// @Param       upload formData []file true "Upload multiple images"
// @Param       product_id query int false "Product ID"
// @Param       is_primary query bool false "Is primary image"
// @Success     200 {object} response.Response{data=[]model.ProductImage} "Successful upload"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/images/upload [post]
func (pr *productRoutes) uploadImages(c *gin.Context) {
	const op = "handlers.product.uploadImages"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	// Получаем ID продукта из query параметра (если есть)
	productIDStr := c.Query("product_id")
	var productID int64
	if productIDStr != "" {
		var err error
		productID, err = strconv.ParseInt(productIDStr, 10, 64)
		if err != nil {
			log.Error("invalid product id", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, response.Error("Invalid product ID"))
			return
		}
	}

	// Получаем флаг основного изображения из query параметра (если есть)
	isPrimaryStr := c.Query("is_primary")
	isPrimary := isPrimaryStr == "true"

	form, err := c.MultipartForm()
	if err != nil {
		log.Error("failed to parse form", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid form data"))
		return
	}

	files, exists := form.File["upload"]
	if !exists || len(files) == 0 {
		log.Warn("no files uploaded")
		c.JSON(http.StatusBadRequest, response.Error("No files uploaded"))
		return
	}

	// Получаем S3 клиент из сервиса
	s3Worker := pr.productService.GetS3Worker()

	// Загружаем файлы в S3
	var uploadedImages []model.ProductImage
	for i, file := range files {
		// Проверяем тип файла
		if !isImageFile(file.Filename) {
			log.Warn("invalid file type", slog.String("filename", file.Filename))
			continue
		}

		// Загружаем файл в S3
		filename, err := s3Worker.UploadFileFromMultipart(file)
		if err != nil {
			log.Error("failed to upload file to S3",
				slog.String("filename", file.Filename),
				slog.String("error", err.Error()),
			)
			continue
		}

		// Получаем URL файла
		fileURL, err := s3Worker.GetFileURL(filename)
		if err != nil {
			log.Error("failed to get file URL",
				slog.String("filename", filename),
				slog.String("error", err.Error()),
			)
			continue
		}

		// Создаем запись об изображении
		image := model.ProductImage{
			ProductID: productID,
			FileUUID:  filename,
			URL:       filename,
			IsPrimary: isPrimary && i == 0, // Только первое изображение может быть основным
		}
		// Устанавливаем временные метки
		image.SetTimestamps()

		// Если указан ID продукта, сохраняем изображение в базе
		if productID > 0 {
			savedImage, err := pr.productService.AddProductImage(c.Request.Context(), image)
			if err != nil {
				log.Error("failed to save image to database",
					slog.String("filename", filename),
					slog.String("error", err.Error()),
				)
				continue
			}
			uploadedImages = append(uploadedImages, *savedImage)
		} else {
			// Если ID продукта не указан, просто добавляем изображение в результат
			uploadedImages = append(uploadedImages, image)
		}

		log.Info("file uploaded successfully",
			slog.String("filename", file.Filename),
			slog.String("s3_filename", filename),
			slog.String("url", fileURL),
		)
	}

	if len(uploadedImages) == 0 {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to upload any images"))
		return
	}

	c.JSON(http.StatusOK, uploadedImages)
}

// isImageFile проверяет, является ли файл изображением по расширению
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}

// getCategories
// @Summary     Get all product categories
// @Description Get all product categories
// @Tags  	    product
// @Produce     json
// @Success     200 {object} response.Response{data=[]model.CategoryFilter} "Successful operation"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/categories [get]
func (pr *productRoutes) getCategories(c *gin.Context) {
	const op = "handlers.product.getCategories"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	// Получаем категории из сервиса
	categories := pr.productService.GetProductCategories()

	log.Info("categories retrieved", slog.Int("count", len(categories)))
	c.JSON(http.StatusOK, categories)
}

// getProduct
// @Summary     Get product by ID
// @Description Get product by ID
// @Tags  	    product
// @Produce     json
// @Param       id path int true "Product ID"
// @Success     200 {object} response.Response{data=model.Product} "Successful operation"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     404 {object} response.Response "Product not found"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/{id} [get]
func (pr *productRoutes) getProduct(c *gin.Context) {
	const op = "handlers.product.getProduct"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid product id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid product ID"))
		return
	}

	// Получение продукта из сервиса
	product, err := pr.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		log.Error("failed to get product", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get product"))
		return
	}

	log.Info("product retrieved", slog.Int64("id", id))
	c.JSON(http.StatusOK, product)
}

// getProductReviews
// @Summary     Get product reviews
// @Description Get product reviews by product ID
// @Tags  	    product
// @Produce     json
// @Param       id path int true "Product ID"
// @Success     200 {object} response.Response{data=[]model.ProductReview} "Successful operation"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     404 {object} response.Response "Product not found"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/{id}/reviews [get]
func (pr *productRoutes) getProductReviews(c *gin.Context) {
	const op = "handlers.product.getProductReviews"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid product id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid product ID"))
		return
	}

	// Получение отзывов из сервиса
	reviews, err := pr.productService.GetProductReviews(c.Request.Context(), id)
	if err != nil {
		log.Error("failed to get product reviews", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get product reviews"))
		return
	}

	log.Info("product reviews retrieved", slog.Int64("product_id", id), slog.Int("count", len(reviews)))
	c.JSON(http.StatusOK, reviews)
}

// addProductReview
// @Summary     Add product review
// @Description Add a new review for a product
// @Tags  	    product
// @Accept      json
// @Produce     json
// @Param       id path int true "Product ID"
// @Param       review body model.ProductReview true "Review data"
// @Success     201 {object} response.Response{data=model.ProductReview} "Review created"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     404 {object} response.Response "Product not found"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/{id}/reviews [post]
func (pr *productRoutes) addProductReview(c *gin.Context) {
	const op = "handlers.product.addProductReview"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid product id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid product ID"))
		return
	}

	var review model.ProductReview
	if err := c.ShouldBindJSON(&review); err != nil {
		log.Error("failed to bind review data", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid review data"))
		return
	}

	// Устанавливаем ID продукта из URL
	review.ProductID = id

	// Устанавливаем текущую дату
	review.Date = time.Now()

	// Добавление отзыва через сервис
	createdReview, err := pr.productService.AddProductReview(c.Request.Context(), review)
	if err != nil {
		log.Error("failed to add product review", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, response.Error("Failed to add product review"))
		return
	}

	review = *createdReview

	log.Info("product review added", slog.Int64("product_id", id), slog.Int64("review_id", review.ID))
	c.JSON(http.StatusCreated, review)
}

// filterProducts
// @Summary     Filter products
// @Description Filter products by various criteria
// @Tags  	    product
// @Accept      json
// @Produce     json
// @Param       filters query model.ProductFilters false "Filter criteria"
// @Success     200 {object} response.Response{data=[]model.Product} "Successful operation"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/filter [get]
func (pr *productRoutes) filterProducts(c *gin.Context) {
	const op = "handlers.product.filterProducts"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var filters model.ProductFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		log.Error("failed to bind filter data", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid filter data"))
		return
	}

	// Фильтрация продуктов через сервис
	products, err := pr.productService.FilterProducts(c.Request.Context(), filters)
	if err != nil {
		log.Error("failed to filter products", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, response.Error("Failed to filter products"))
		return
	}

	log.Info("products filtered", slog.Int("count", len(products)))
	c.JSON(http.StatusOK, products)
}

// createProduct
// @Summary     Create a new product
// @Description Create a new product
// @Tags  	    product
// @Accept      json
// @Produce     json
// @Param       product body model.ProductCreateRequest true "Product data"
// @Success     201 {object} response.Response{data=model.Product} "Product created"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product [post]
func (pr *productRoutes) createProduct(c *gin.Context) {
	const op = "handlers.product.createProduct"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var productRequest model.ProductCreateRequest
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		log.Error("failed to bind product data", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid product data: "+err.Error()))
		return
	}

	// Преобразуем запрос в модель продукта
	product := productRequest.ToProduct()

	// Создание продукта через сервис
	createdProduct, err := pr.productService.CreateProduct(c.Request.Context(), product)
	if err != nil {
		log.Error("failed to create product", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, response.Error("Failed to create product: "+err.Error()))
		return
	}

	log.Info("product created", slog.Int64("id", createdProduct.ID))
	c.JSON(http.StatusCreated, createdProduct)
}

// updateProduct
// @Summary     Update a product
// @Description Update an existing product
// @Tags  	    product
// @Accept      json
// @Produce     json
// @Param       id path int true "Product ID"
// @Param       product body model.ProductUpdateRequest true "Product data"
// @Success     200 {object} response.Response{data=model.Product} "Product updated"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     404 {object} response.Response "Product not found"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/{id} [put]
func (pr *productRoutes) updateProduct(c *gin.Context) {
	const op = "handlers.product.updateProduct"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid product id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid product ID"))
		return
	}

	// Получаем существующий продукт
	existingProduct, err := pr.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		log.Error("failed to get product", slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, response.Error("Product not found"))
		return
	}

	var updateRequest model.ProductUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Error("failed to bind product data", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid product data: "+err.Error()))
		return
	}

	// Применяем изменения к существующему продукту
	updateRequest.ApplyToProduct(existingProduct)

	// Обновление продукта через сервис
	updatedProduct, err := pr.productService.UpdateProduct(c.Request.Context(), *existingProduct)
	if err != nil {
		log.Error("failed to update product", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, response.Error("Failed to update product: "+err.Error()))
		return
	}

	log.Info("product updated", slog.Int64("id", updatedProduct.ID))
	c.JSON(http.StatusOK, updatedProduct)
}

// deleteProduct
// @Summary     Delete a product
// @Description Delete a product by ID
// @Tags  	    product
// @Produce     json
// @Param       id path int true "Product ID"
// @Success     200 {object} response.Response "Product deleted"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     404 {object} response.Response "Product not found"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/{id} [delete]
func (pr *productRoutes) deleteProduct(c *gin.Context) {
	const op = "handlers.product.deleteProduct"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid product id", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid product ID"))
		return
	}

	// Удаление продукта через сервис
	if err := pr.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		log.Error("failed to delete product", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, response.Error("Failed to delete product"))
		return
	}

	log.Info("product deleted", slog.Int64("id", id))
	c.JSON(http.StatusNoContent, nil)
}

// UploadReviewFile
// @Summary     Upload multiple images for review
// @Description Upload multiple images for review
// @Tags  	    product
// @Accept      multipart/form-data
// @Produce     json
// @Param       upload formData []file true "Upload multiple images"
// @Param       review_id query int false "Review ID"
// @Param       is_primary query bool false "Is primary image"
// @Success     200 {object} response.Response{data=[]model.ProductImage} "Successful upload"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/{id}/images/upload [post]
func (pr *productRoutes) UploadReviewFile(c *gin.Context) {
	const op = "handlers.product.UploadReviewFile"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	// Получаем ID продукта из query параметра (если есть)
	reviewIDStr := c.Query("review_id")
	var reviewID int64
	if reviewIDStr != "" {
		var err error
		reviewID, err = strconv.ParseInt(reviewIDStr, 10, 64)
		if err != nil {
			log.Error("invalid review id", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, response.Error("Invalid review ID"))
			return
		}
	}

	// Получаем флаг основного изображения из query параметра (если есть)
	isPrimaryStr := c.Query("is_primary")
	isPrimary := isPrimaryStr == "true"

	form, err := c.MultipartForm()
	if err != nil {
		log.Error("failed to parse form", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Invalid form data"))
		return
	}

	files, exists := form.File["upload"]
	if !exists || len(files) == 0 {
		log.Warn("no files uploaded")
		c.JSON(http.StatusBadRequest, response.Error("No files uploaded"))
		return
	}

	// Получаем S3 клиент из сервиса
	s3Worker := pr.productService.GetS3WorkerReview()

	// Загружаем файлы в S3
	var uploadedImages []model.ReviewImages
	for i, file := range files {
		// Проверяем тип файла
		if !isImageFile(file.Filename) {
			log.Warn("invalid file type", slog.String("filename", file.Filename))
			continue
		}

		// Загружаем файл в S3
		filename, err := s3Worker.UploadFileFromMultipart(file)
		if err != nil {
			log.Error("failed to upload file to S3",
				slog.String("filename", file.Filename),
				slog.String("error", err.Error()),
			)
			continue
		}

		// Получаем URL файла
		fileURL, err := s3Worker.GetFileURL(filename)
		if err != nil {
			log.Error("failed to get file URL",
				slog.String("filename", filename),
				slog.String("error", err.Error()),
			)
			continue
		}

		// Создаем запись об изображении
		image := model.ReviewImages{
			ReviewID:  reviewID,
			FileUUID:  filename,
			URL:       filename,
			IsPrimary: isPrimary && i == 0, // Только первое изображение может быть основным
		}
		// Устанавливаем временные метки
		image.SetTimestamps()

		// Если указан ID продукта, сохраняем изображение в базе
		if reviewID > 0 {
			savedImage, err := pr.productService.AddReviewImage(image)
			if err != nil {
				log.Error("failed to save image to database",
					slog.String("filename", filename),
					slog.String("error", err.Error()),
				)
				continue
			}
			uploadedImages = append(uploadedImages, savedImage)
		} else {
			// Если ID продукта не указан, просто добавляем изображение в результат
			uploadedImages = append(uploadedImages, image)
		}

		log.Info("file uploaded successfully",
			slog.String("filename", file.Filename),
			slog.String("s3_filename", filename),
			slog.String("url", fileURL),
		)
	}

	if len(uploadedImages) == 0 {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to upload any images"))
		return
	}

	c.JSON(http.StatusOK, uploadedImages)
}
