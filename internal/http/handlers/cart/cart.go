package cart

import (
	"github.com/RCSE2025/backend-go/internal/http/middleware/auth"
	"github.com/RCSE2025/backend-go/internal/http/middleware/logger"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/RCSE2025/backend-go/pkg/logger/sl"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type cartRoutes struct {
	cartService *service.CartService
}

func NewCartRoutes(h *gin.RouterGroup, cs *service.CartService, jwtService service.JWTService) {
	g := h.Group("/cart")

	ur := cartRoutes{cartService: cs}

	validateJWTmw := auth.ValidateJWT(jwtService)

	g.GET("", validateJWTmw, ur.GetCartProduct)
	g.POST("", validateJWTmw, ur.PostCartProduct)
	g.DELETE("", validateJWTmw, ur.DeleteCartProduct)
	g.PUT("", validateJWTmw, ur.SetCartQuantity)
}

type PostProductRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// PostCartProduct
// @Summary 	Post product in cart
// @Description Post product in cart
// @Tags		cart
// @Accept		json
// @Produce		json
// @Failure 	500 {object} response.Response
// @Failure		400 {object} response.Response
// @Success 	200 {object} model.CartItem
// @Router		/cart [post]
// @Param  request body PostProductRequest true "request"
// @Security OAuth2PasswordBearer
func (cr *cartRoutes) PostCartProduct(c *gin.Context) {
	const op = "handlers.cart.PostCartProduct"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)))

	var req PostProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("cannot bind request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	cart, err := cr.cartService.PostInCart(c.GetInt64("user_id"), req.ProductID, req.Quantity)
	if err != nil {
		log.Error("cannot post product in cart", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, cart)
}

type ProductsID []int64

// DeleteCartProduct
// @Summary 	Delete product in cart
// @Description Delete product in user cart with ids
// @Tags		cart
// @Accept		json
// @Produce		json
// @Failure 	500 {object} response.Response
// @Success		200 {object} response.Response
// @Router		/cart [delete]
// @Param  request body ProductsID true "request"
// @Security OAuth2PasswordBearer
func (cr *cartRoutes) DeleteCartProduct(c *gin.Context) {
	const op = "handlers.cart.DeleteCartProduct"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)))

	var req ProductsID
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("cannot bind request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	err := cr.cartService.DeleteCart(c.GetInt64("user_id"), req)
	if err != nil {
		log.Error("cannot delete product in user cart", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.OK())
}

// GetCartProduct
// @Summary 	Get user cart
// @Description Get user cart
// @Tags		cart
// @Accept		json
// @Produce		json
// @Failure 	500 {object} response.Response
// @Success		200 {object} []model.CartItemsResponse
// @Router		/cart [get]
// @Security OAuth2PasswordBearer
func (cr *cartRoutes) GetCartProduct(c *gin.Context) {
	const op = "handlers.cart.GetCartProduct"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)))

	userCard, err := cr.cartService.GetUserCart(c.GetInt64("user_id"))
	if err != nil {
		log.Error("cannot get product in user cart", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	if userCard == nil {
		userCard = make([]model.CartItemsResponse, 0)
	}
	c.JSON(http.StatusOK, userCard)
}

type SetQuantityRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// SetCartQuantity
// @Summary 	Set product quantity in cart
// @Description Set product quantity in cart
// @Tags		cart
// @Accept		json
// @Produce		json
// @Failure 	500 {object} response.Response
// @Success		200 {object} response.Response
// @Router		/cart [put]
// @Param request body map[int64]int true "request"
// @Security OAuth2PasswordBearer
func (cr *cartRoutes) SetCartQuantity(c *gin.Context) {
	const op = "handlers.cart.SetCartQuantity"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)))

	var req map[int64]int
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("cannot bind request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	for id, quantity := range req {
		err := cr.cartService.SetCartQuantity(c.GetInt64("user_id"), id, quantity)
		if err != nil {
			log.Error("cannot set product quantity in user cart", sl.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
			return
		}
	}
	c.JSON(http.StatusOK, response.OK())
}
