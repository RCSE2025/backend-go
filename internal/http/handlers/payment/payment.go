package payment

import (
	"fmt"
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type paymentRoutes struct {
	yookassa *service.YookassaPayment
}

func NewProductRoutes(h *gin.RouterGroup, yookassa *service.YookassaPayment) {
	g := h.Group("/payment")

	pr := paymentRoutes{
		yookassa: yookassa,
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
// @Tags Payment
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
	tt := t["metadata"].(map[string]interface{})["order_id"]
	fmt.Println(tt)
	fmt.Println(tt.(string))
	fmt.Println(tt.(int))
	var req notification
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	fmt.Println(req)
}
