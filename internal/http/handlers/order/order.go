package order

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

type orderRoutes struct {
	ordService  *service.OrderService
	cartService *service.CartService
}

func NewOrderRoutes(h *gin.RouterGroup, s *service.OrderService, jwtService service.JWTService, cartService *service.CartService) {
	g := h.Group("/order")

	validateJWTmw := auth.ValidateJWT(jwtService)
	ordR := orderRoutes{ordService: s, cartService: cartService}

	g.POST("/create_order_manual", validateJWTmw, ordR.CreateOrder)
	g.POST("/create_order_yookassa", validateJWTmw, ordR.CreateOrderYookassa)
	g.GET("", validateJWTmw, ordR.GetListOrders)
	g.PUT("", validateJWTmw, ordR.SetOrderStatus)
}

type CreateOrderRequest struct {
	Address string `json:"address"`
}

// CreateOrder
// @Summary 	Create Order
// @Description Create new order
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Param request body CreateOrderRequest true "request"
// @Success     201 {object} response.Response`
// @Router      /order/create_order_manual [post]
// @Security OAuth2PasswordBearer
func (ordR *orderRoutes) CreateOrder(c *gin.Context) {
	const op = "handlers.user.CreateOrder"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var req CreateOrderRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error("cannot parse request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userID := c.GetInt64("user_id")
	order, err := ordR.ordService.CreateOrder(userID, req.Address)
	if err != nil {
		log.Error("can't create order", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Can't create order"))
		return
	}

	var deleteItemsIDs []int64
	cart, err := ordR.cartService.GetUserCart(userID)
	if err != nil {
		log.Error("can't get user cart", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Can't get user cart"))
		return
	}
	for _, r := range cart {
		_, err := ordR.ordService.CreateOrderItem(userID, order.ID, r.ProductID, r.Quantity)
		if err != nil {
			log.Error("can't create order items", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, response.Error("Can't create order items"))
			return
		}
		deleteItemsIDs = append(deleteItemsIDs, r.ProductID)
	}

	err = ordR.cartService.DeleteCart(userID, deleteItemsIDs)
	if err != nil {
		log.Error("can't delete cart items", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Can't create cart items"))
		return
	}
	err = ordR.ordService.ConfirmOrderPayment(order.ID)
	if err != nil {
		log.Error("can't confirm order payment", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("can't confirm order payment"))
		return
	}
	c.JSON(http.StatusCreated, response.OK())
}

// GetListOrders
// @Summary 	Get List Order
// @Description Get List Order
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Success     200 {object} []model.OrderItemResponse
// @Router      /order [get]
// @Security OAuth2PasswordBearer
func (ordR *orderRoutes) GetListOrders(c *gin.Context) {
	const op = "handlers.user.GetListOrders"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	userID := c.GetInt64("user_id")
	orders, err := ordR.ordService.GetUserOrders(userID)
	if err != nil {
		log.Error("can't get user orders", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("can't get user orders"))
		return
	}
	c.JSON(http.StatusOK, orders)

}

type SetOrderStatusRequest struct {
	OrderID int64                 `json:"order_id"`
	Status  model.OrderStatusType `json:"status"`
}

// SetOrderStatus
// @Summary 	Set Order status
// @Description Set Order status
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Success     200 {object} response.Response`
// @Router      /order [put]
// @Param request body SetOrderStatusRequest true "request"
// @Security OAuth2PasswordBearer
func (ordR *orderRoutes) SetOrderStatus(c *gin.Context) {
	const op = "handlers.user.SetOrderStatus"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var req SetOrderStatusRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error("cannot parse request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userID := c.GetInt64("user_id")
	err := ordR.ordService.SetOrderStatus(userID, req.OrderID, req.Status)
	if err != nil {
		log.Error("can't set order status", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("can't set order status"))
		return
	}
	c.JSON(http.StatusOK, response.OK())
}

type CreateOrderYookassaRequest struct {
	ConfirmURL string `json:"confirm_url"`
}

// CreateOrderYookassa
// @Summary 	Create Order
// @Description Create new order
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Failure     500 {object} response.Response
// @Param request body CreateOrderRequest true "request"
// @Success     201 {object} CreateOrderYookassaRequest`
// @Router      /order/create_order_yookassa [post]
// @Security OAuth2PasswordBearer
func (ordR *orderRoutes) CreateOrderYookassa(c *gin.Context) {
	const op = "handlers.user.CreateOrder"
	log := logger.FromContext(c).With(
		slog.String("op", op),
		slog.String("request_id", requestid.Get(c)),
	)

	var req CreateOrderRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error("cannot parse request", sl.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userID := c.GetInt64("user_id")
	order, err := ordR.ordService.CreateOrder(userID, req.Address)
	if err != nil {
		log.Error("can't create order", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Can't create order"))
		return
	}

	var deleteItemsIDs []int64
	cart, err := ordR.cartService.GetUserCart(userID)
	if err != nil {
		log.Error("can't get user cart", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Can't get user cart"))
		return
	}
	var totalPrice float64

	for _, r := range cart {
		_, err := ordR.ordService.CreateOrderItem(userID, order.ID, r.ProductID, r.Quantity)
		if err != nil {
			log.Error("can't create order items", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, response.Error("Can't create order items"))
			return
		}

		totalPrice += float64(r.Quantity) * r.Product.Price
		deleteItemsIDs = append(deleteItemsIDs, r.ProductID)
	}

	err = ordR.cartService.DeleteCart(userID, deleteItemsIDs)
	if err != nil {
		log.Error("can't delete cart items", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Can't create cart items"))
		return
	}

	url, err := ordR.ordService.CreateOrderPayment(order.ID, totalPrice)
	if err != nil {
		log.Error("can't create order payment", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, response.Error("Can't create order payment"))
		return
	}

	c.JSON(http.StatusOK, CreateOrderYookassaRequest{ConfirmURL: url})
}
