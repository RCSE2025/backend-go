package service

import (
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	"github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type YookassaPayment struct {
	client *yookassa.Client
}

func NewYookassaPayment(accountId string, secretKey string) *YookassaPayment {
	return &YookassaPayment{
		client: yookassa.NewClient(accountId, secretKey),
	}
}

func (p *YookassaPayment) CreatePayment(paymentData map[string]interface{}) (string, error) {
	paymentHandler := yookassa.NewPaymentHandler(p.client)
	// Создаем платеж
	payment, _ := paymentHandler.CreatePayment(&yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    "1000.00",
			Currency: "RUB",
		},
		PaymentMethod: yoopayment.PaymentMethodType("bank_card"),
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: "https://www.example.com",
		},
		Description: "Test payment",
	})

	return payment.Confirmation.(map[string]interface{})["confirmation_url"].(string), nil
}
