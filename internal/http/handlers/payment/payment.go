package payment

import (
	"fmt"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type paymentRoutes struct {
	yookassa     *service.YookassaPayment
	orderService *service.OrderService
}

func NewProductRoutes(h *gin.RouterGroup, yookassa *service.YookassaPayment, orderService *service.OrderService) {
	g := h.Group("/payment")

	pr := paymentRoutes{
		yookassa:     yookassa,
		orderService: orderService,
	}

	g.POST("/notifications", pr.notification)
}

type orderInfo struct {
	OrderId string `json:"order_id"`
}

type notification struct {
	Event    string    `json:"event"`
	Metadata orderInfo `json:"metadata"`
}

// @Summary Payment notification
// @Tags payment
// @Accept json
// @Produce json
// @Success 200
// @Router /payment/notifications [post]
func (pr *paymentRoutes) notification(c *gin.Context) {
	var t map[string]interface{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	fmt.Println(t)

	if object, ok := t["object"].(map[string]interface{}); ok {
		if metadata, ok := object["metadata"].(map[string]interface{}); ok {
			if orderID, ok := metadata["order_id"].(string); ok {
				fmt.Println("Order ID:", orderID)
				i, err := strconv.ParseInt(orderID, 10, 64)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
					return
				}
				err = pr.orderService.ConfirmOrderPayment(i)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
					return
				}
				c.Status(http.StatusOK)
				return
			}
		}
	}
}
