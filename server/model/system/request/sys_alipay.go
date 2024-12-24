package request

// AlipayPaymentRequest represents the request to create an Alipay payment
type AlipayPaymentRequest struct {
	OrderId string  `json:"orderId" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Subject string  `json:"subject" binding:"required"`
}
