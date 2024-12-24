package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type PaymentRouter struct{}

func (s *PaymentRouter) InitPaymentRouter(Router *gin.RouterGroup) {
	paymentRouter := Router.Group("payment")
	alipayApi := v1.ApiGroupApp.SystemApiGroup.AlipayApi
	{
		// Alipay endpoints
		paymentRouter.POST("alipay/create", alipayApi.CreateAlipayPayment)
		paymentRouter.GET("alipay/status", alipayApi.GetAlipayStatus)
	}
}
