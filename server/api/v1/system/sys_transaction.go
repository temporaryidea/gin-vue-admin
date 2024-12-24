package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TransactionApi struct{}

// CreateTransaction godoc
// @Tags      Transaction
// @Summary   Create a new transaction
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.CreateTransactionRequest  true  "Transaction information"
// @Success   200   {object}  response.Response{msg=string}      "Create transaction"
// @Router    /transaction/create [post]
func (t *TransactionApi) CreateTransaction(c *gin.Context) {
	var req systemReq.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Create transaction logic here
	transaction := &system.SysTransaction{
		ProductID:   req.ProductID,
		Amount:      req.Amount,
		PaymentType: req.PaymentType,
		Description: req.Description,
		Status:      "pending",
	}

	if err := global.GVA_DB.Create(transaction).Error; err != nil {
		global.GVA_LOG.Error("Create transaction failed!", zap.Error(err))
		response.FailWithMessage("Create transaction failed", c)
		return
	}

	response.OkWithMessage("Transaction created successfully", c)
}

// GetTransactionList godoc
// @Tags      Transaction
// @Summary   Get transaction list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.GetTransactionListRequest                        true  "Transaction list parameters"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}    "Get transaction list"
// @Router    /transaction/list [post]
func (t *TransactionApi) GetTransactionList(c *gin.Context) {
	var req systemReq.GetTransactionListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var list []system.SysTransaction
	var total int64

	db := global.GVA_DB.Model(&system.SysTransaction{})

	// Apply filters
	if req.OrderID != "" {
		db = db.Where("order_id = ?", req.OrderID)
	}
	if req.UserID != 0 {
		db = db.Where("user_id = ?", req.UserID)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.PaymentType != "" {
		db = db.Where("payment_type = ?", req.PaymentType)
	}

	err := db.Count(&total).Error
	if err != nil {
		global.GVA_LOG.Error("Get transaction count failed!", zap.Error(err))
		response.FailWithMessage("Get transaction list failed", c)
		return
	}

	err = db.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Find(&list).Error
	if err != nil {
		global.GVA_LOG.Error("Get transaction list failed!", zap.Error(err))
		response.FailWithMessage("Get transaction list failed", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "Get transaction list successfully", c)
}

// UpdateTransactionStatus godoc
// @Tags      Transaction
// @Summary   Update transaction status
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.UpdateTransactionStatusRequest  true  "Transaction status update"
// @Success   200   {object}  response.Response{msg=string}           "Update transaction status"
// @Router    /transaction/status [post]
func (t *TransactionApi) UpdateTransactionStatus(c *gin.Context) {
	var req systemReq.UpdateTransactionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := global.GVA_DB.Model(&system.SysTransaction{}).
		Where("order_id = ?", req.OrderID).
		Update("status", req.Status).Error

	if err != nil {
		global.GVA_LOG.Error("Update transaction status failed!", zap.Error(err))
		response.FailWithMessage("Update transaction status failed", c)
		return
	}

	response.OkWithMessage("Transaction status updated successfully", c)
}
