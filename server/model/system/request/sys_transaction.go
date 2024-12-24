package request

// CreateTransactionRequest represents the request structure for creating a transaction
type CreateTransactionRequest struct {
	ProductID   uint    `json:"productId" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	PaymentType string  `json:"paymentType" binding:"required"`
	Description string  `json:"description"`
}

// GetTransactionListRequest represents the request structure for getting transaction list
type GetTransactionListRequest struct {
	OrderID     string `json:"orderId"`
	UserID      uint   `json:"userId"`
	Status      string `json:"status"`
	PaymentType string `json:"paymentType"`
	PageInfo
}

// UpdateTransactionStatusRequest represents the request structure for updating transaction status
type UpdateTransactionStatusRequest struct {
	OrderID string `json:"orderId" binding:"required"`
	Status  string `json:"status" binding:"required"`
}
