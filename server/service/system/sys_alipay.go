package system

import (
	"errors"
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/viper"
)

type AlipayService struct {
	Client *alipay.Client
}

var AlipayServiceApp = new(AlipayService)

// NewAlipayService initializes a new AlipayService with configuration from alipay.yaml
func NewAlipayService() (*AlipayService, error) {
	viper.SetConfigName("alipay")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read alipay config: %v", err)
	}

	// Load configuration
	appID := viper.GetString("alipay.app_id")
	privateKey := viper.GetString("alipay.private_key")
	publicKey := viper.GetString("alipay.public_key")
	sandboxMode := viper.GetBool("alipay.sandbox_mode")

	if appID == "" || privateKey == "" || publicKey == "" {
		return nil, errors.New("missing required Alipay configuration")
	}

	// Initialize Alipay client
	client, err := alipay.New(appID, privateKey, sandboxMode)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Alipay client: %v", err)
	}

	// Load Alipay public key
	err = client.LoadAliPayPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load Alipay public key: %v", err)
	}

	return &AlipayService{Client: client}, nil
}

// GenerateQRCode creates a new Alipay trade and returns the QR code URL
func (s *AlipayService) GenerateQRCode(orderId string, amount float64, subject string) (string, error) {
	if s.Client == nil {
		return "", errors.New("Alipay client not initialized")
	}

	// Create trade precreate request
	p := alipay.TradePreCreate{
		Trade: alipay.Trade{
			Subject:     subject,
			OutTradeNo: orderId,
			TotalAmount: fmt.Sprintf("%.2f", amount),
		},
	}

	// Execute the request
	rsp, err := s.Client.TradePreCreate(p)
	if err != nil {
		return "", fmt.Errorf("failed to create Alipay trade: %v", err)
	}

	if !rsp.IsSuccess() {
		return "", fmt.Errorf("Alipay trade creation failed: %s", rsp.SubMsg)
	}

	return rsp.QRCode, nil
}

// QueryPaymentStatus checks the status of a payment by order ID
func (s *AlipayService) QueryPaymentStatus(orderId string) (string, error) {
	if s.Client == nil {
		return "", errors.New("Alipay client not initialized")
	}

	// Create trade query request
	p := alipay.TradeQuery{
		OutTradeNo: orderId,
	}

	// Execute the request
	rsp, err := s.Client.TradeQuery(p)
	if err != nil {
		return "", fmt.Errorf("failed to query Alipay trade: %v", err)
	}

	if !rsp.IsSuccess() {
		return "", fmt.Errorf("Alipay trade query failed: %s", rsp.SubMsg)
	}


	return rsp.TradeStatus, nil
}

// Initialize initializes the AlipayService singleton
func (s *AlipayService) Initialize() error {
	service, err := NewAlipayService()
	if err != nil {
		return err
	}
	*s = *service
	return nil
}
