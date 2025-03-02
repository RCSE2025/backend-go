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
	//fmt.Println(t)
	//fmt.Println(t["metadata"])
	//tt := t["metadata"].(map[string]string)
	//fmt.Println(tt)
	//fmt.Println(tt["order_id"])
	////fmt.Println(tt.(int))
	//var req notification
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))
	//	return
	//}
	//fmt.Println(req)

	fmt.Println(t)

	// The metadata is nested inside the "object" field, need to safely access it
	if object, ok := t["object"].(map[string]interface{}); ok {
		if metadata, ok := object["metadata"].(map[string]interface{}); ok {
			if orderID, ok := metadata["order_id"].(string); ok {
				fmt.Println("Order ID:", orderID)
				// Process the order ID

				// Return successful response
				c.Status(http.StatusOK)
				return
			}
		}
	}
}
