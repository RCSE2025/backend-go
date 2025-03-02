package service

import (
	"fmt"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	"github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
	"math"
)

type YookassaPayment struct {
	client *yookassa.Client
}

func NewYookassaPayment(accountId string, secretKey string) *YookassaPayment {
	return &YookassaPayment{
		client: yookassa.NewClient(accountId, secretKey),
	}
}

func (p *YookassaPayment) CreateOrderPayment(orderID int64, amount float64) (string, error) {
	paymentHandler := yookassa.NewPaymentHandler(p.client)
	// Создаем платеж
	payment, _ := paymentHandler.CreatePayment(&yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    fmt.Sprintf("%.2f", math.Round(amount*100)/100),
			Currency: "RUB",
		},
		PaymentMethod: yoopayment.PaymentMethodType("bank_card"),
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: "https://www.example.com",
		},
		Description: "Test payment",
		Metadata: map[string]interface{}{
			"order_id": orderID,
		},
	})

	return payment.Confirmation.(map[string]interface{})["confirmation_url"].(string), nil
}
