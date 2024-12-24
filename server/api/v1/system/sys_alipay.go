package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/gin-gonic/gin"
)

type AlipayApi struct{}

// CreateAlipayPayment generates a QR code for Alipay payment
func (a *AlipayApi) CreateAlipayPayment(c *gin.Context) {
	var req request.AlipayPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Invalid request parameters", c)
		return
	}

	qrCode, err := system.AlipayServiceApp.GenerateQRCode(req.OrderId, req.Amount, req.Subject)
	if err != nil {
		global.GVA_LOG.Error("Failed to generate Alipay QR code:", err)
		response.FailWithMessage("Failed to generate payment QR code", c)
		return
	}

	response.OkWithData(gin.H{
		"qrCode": qrCode,
		"orderId": req.OrderId,
	}, c)
}

// GetAlipayStatus checks the payment status
func (a *AlipayApi) GetAlipayStatus(c *gin.Context) {
	orderId := c.Query("orderId")
	if orderId == "" {
		response.FailWithMessage("Order ID is required", c)
		return
	}

	status, err := system.AlipayServiceApp.QueryPaymentStatus(orderId)
	if err != nil {
		global.GVA_LOG.Error("Failed to query Alipay payment status:", err)
		response.FailWithMessage("Failed to check payment status", c)
		return
	}

	response.OkWithData(gin.H{
		"orderId": orderId,
		"status": status,
	}, c)
}
