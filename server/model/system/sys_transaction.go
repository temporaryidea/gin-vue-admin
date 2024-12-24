package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SysTransaction represents a transaction record in the system
type SysTransaction struct {
	global.GVA_MODEL
	OrderID      string  `json:"orderId" gorm:"uniqueIndex;comment:订单ID"`
	UserID       uint    `json:"userId" gorm:"comment:用户ID"`
	ProductID    uint    `json:"productId" gorm:"comment:产品ID"`
	Amount       float64 `json:"amount" gorm:"comment:交易金额"`
	Status       string  `json:"status" gorm:"comment:交易状态"` // pending, completed, failed
	PaymentType  string  `json:"paymentType" gorm:"comment:支付方式"`
	Description  string  `json:"description" gorm:"comment:交易描述"`
}

func (SysTransaction) TableName() string {
	return "sys_transactions"
}

// SysProduct represents a product in the system
type SysProduct struct {
	global.GVA_MODEL
	Name        string  `json:"name" gorm:"comment:产品名称"`
	Description string  `json:"description" gorm:"comment:产品描述"`
	Price       float64 `json:"price" gorm:"comment:产品价格"`
	Stock       int     `json:"stock" gorm:"comment:库存数量"`
	Category    string  `json:"category" gorm:"comment:产品类别"`
	ImageURL    string  `json:"imageUrl" gorm:"comment:产品图片URL"`
}

func (SysProduct) TableName() string {
	return "sys_products"
}
