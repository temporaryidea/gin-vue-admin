package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PaymentRouter struct{}

func (s *PaymentRouter) InitPaymentRouter(Router *gin.RouterGroup) {
	paymentRouter := Router.Group("payment").Use(middleware.OperationRecord())
	paymentRouterWithoutRecord := Router.Group("payment")
	var paymentApi = v1.ApiGroupApp.SystemApiGroup.TransactionApi
	{
		paymentRouter.POST("create", paymentApi.CreateTransaction)   // Create new transaction
		paymentRouter.POST("refund", paymentApi.RefundTransaction)  // Process refund
	}
	{
		paymentRouterWithoutRecord.GET("status", paymentApi.GetTransactionStatus) // Get transaction status
	}
}
