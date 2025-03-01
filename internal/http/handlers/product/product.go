package product

import (
	"fmt"
	"github.com/RCSE2025/backend-go/internal/http/middleware/logger"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type productRoutes struct {
	//s *service.UserService
}

func NewProductRoutes(h *gin.RouterGroup, jwtService service.JWTService) {
	g := h.Group("/product")

	ur := productRoutes{}

	g.POST("/images/upload", ur.uploadImages)
}

// uploadImages
// @Summary     Upload multiple images
// @Description Upload multiple images
// @Tags  	    product
// @Accept      multipart/form-data
// @Produce     json
// @Param       upload formData []file true "Upload multiple images"
// @Success     200 {object} response.Response "Successful upload"
// @Failure     400 {object} response.Response "Bad request"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /product/images/upload [post]
func (pr *productRoutes) uploadImages(c *gin.Context) {
	const op = "handlers.product.uploadImages"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	form, err := c.MultipartForm()
	if err != nil {
		log.Error("failed to parse form", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Response{Message: "Invalid form data"})
		return
	}

	files, exists := form.File["upload"]
	if !exists || len(files) == 0 {
		log.Warn("no files uploaded")
		c.JSON(http.StatusBadRequest, response.Response{Message: "No files uploaded"})
		return
	}

	for _, file := range files {
		log.Info("file uploaded", slog.String("filename", file.Filename))
		// c.SaveUploadedFile(file, dst) // Раскомментируйте, если нужно сохранять файлы
	}

	c.JSON(http.StatusOK, response.Response{Message: fmt.Sprintf("%d files uploaded!", len(files))})
}
