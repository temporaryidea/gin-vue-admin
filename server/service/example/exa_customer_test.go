package example

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
)

func TestCreateExaCustomer(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Auto migrate customer table
	if err := global.GVA_DB.AutoMigrate(&example.ExaCustomer{}); err != nil {
		t.Fatalf("Failed to migrate customer table: %v", err)
	}

	tests := []struct {
		name    string
		e       example.ExaCustomer
		wantErr bool
	}{
		{
			name: "Test Create Valid Customer",
			e: example.ExaCustomer{
				CustomerName:       "Test Customer",
				CustomerPhoneData: "1234567890",
				SysUserID:         1,
				SysUserAuthorityID: 1,
			},
			wantErr: false,
		},
		{
			name: "Test Create Customer Without Name",
			e: example.ExaCustomer{
				CustomerPhoneData: "1234567890",
				SysUserID:         1,
				SysUserAuthorityID: 1,
			},
			wantErr: true,
		},
		{
			name: "Test Create Customer Without Phone",
			e: example.ExaCustomer{
				CustomerName:       "Test Customer",
				SysUserID:         1,
				SysUserAuthorityID: 1,
			},
			wantErr: true,
		},
		{
			name: "Test Create Customer Without UserID",
			e: example.ExaCustomer{
				CustomerName:       "Test Customer",
				CustomerPhoneData: "1234567890",
				SysUserAuthorityID: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CustomerService{}
			err := service.CreateExaCustomer(tt.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateExaCustomer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteExaCustomer(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Auto migrate customer table
	if err := global.GVA_DB.AutoMigrate(&example.ExaCustomer{}); err != nil {
		t.Fatalf("Failed to migrate customer table: %v", err)
	}

	// Create test customer for deletion
	testCustomer := example.ExaCustomer{
		CustomerName:       "Test Customer",
		CustomerPhoneData: "1234567890",
		SysUserID:         1,
		SysUserAuthorityID: 1,
	}
	if err := global.GVA_DB.Create(&testCustomer).Error; err != nil {
		t.Fatalf("Failed to create test customer: %v", err)
	}

	tests := []struct {
		name    string
		e       example.ExaCustomer
		wantErr bool
	}{
		{
			name: "Test Delete Existing Customer",
			e: example.ExaCustomer{
				CustomerName:       "Test Customer",
				CustomerPhoneData: "1234567890",
				SysUserID:         1,
				SysUserAuthorityID: 1,
			},
			wantErr: false,
		},
		{
			name: "Test Delete Non-existent Customer",
			e: example.ExaCustomer{
				SysUserID:         999,
				SysUserAuthorityID: 1,
			},
			wantErr: true,
		},
		{
			name: "Test Delete Customer Without ID",
			e: example.ExaCustomer{
				CustomerName:       "Test Customer",
				CustomerPhoneData: "1234567890",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CustomerService{}
			err := service.DeleteExaCustomer(tt.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteExaCustomer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TestUpdateExaCustomer(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Auto migrate customer table
	if err := global.GVA_DB.AutoMigrate(&example.ExaCustomer{}); err != nil {
		t.Fatalf("Failed to migrate customer table: %v", err)
	}

	// Create test customer for update
	testCustomer := example.ExaCustomer{
		CustomerName:       "Original Name",
		CustomerPhoneData: "1234567890",
		SysUserID:         1,
		SysUserAuthorityID: 1,
	}
	if err := global.GVA_DB.Create(&testCustomer).Error; err != nil {
		t.Fatalf("Failed to create test customer: %v", err)
	}

	tests := []struct {
		name    string
		e       *example.ExaCustomer
		wantErr bool
	}{
		{
			name: "Test Update Valid Customer",
			e: &example.ExaCustomer{
				CustomerName:       "Updated Customer",
				CustomerPhoneData: "9876543210",
				SysUserID:         1,
				SysUserAuthorityID: 1,
			},
			wantErr: false,
		},
		{
			name: "Test Update Non-existent Customer",
			e: &example.ExaCustomer{
				SysUserID:         999,
				SysUserAuthorityID: 1,
			},
			wantErr: true,
		},
		{
			name:    "Test Update Nil Customer",
			e:       nil,
			wantErr: true,
		},
		{
			name: "Test Update Customer Without Required Fields",
			e: &example.ExaCustomer{
				CustomerName: "Only Name",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CustomerService{}
			err := service.UpdateExaCustomer(tt.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateExaCustomer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetExaCustomer(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Auto migrate customer table
	if err := global.GVA_DB.AutoMigrate(&example.ExaCustomer{}); err != nil {
		t.Fatalf("Failed to migrate customer table: %v", err)
	}

	// Create test customer for retrieval
	testCustomer := example.ExaCustomer{
		CustomerName:       "Test Customer",
		CustomerPhoneData: "1234567890",
		SysUserID:         1,
		SysUserAuthorityID: 1,
	}
	if err := global.GVA_DB.Create(&testCustomer).Error; err != nil {
		t.Fatalf("Failed to create test customer: %v", err)
	}

	tests := []struct {
		name       string
		id         uint
		wantErr    bool
		checkProps bool
	}{
		{
			name:       "Test Get Existing Customer",
			id:         1,
			wantErr:    false,
			checkProps: true,
		},
		{
			name:       "Test Get Non-existent Customer",
			id:         999,
			wantErr:    true,
			checkProps: false,
		},
		{
			name:       "Test Get Customer With Zero ID",
			id:         0,
			wantErr:    true,
			checkProps: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CustomerService{}
			customer, err := service.GetExaCustomer(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExaCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkProps {
				if customer.ID != tt.id {
					t.Errorf("GetExaCustomer() customer.ID = %v, want %v", customer.ID, tt.id)
				}
				// Add more property checks as needed
			}
		})
	}
}

func TestGetCustomerInfoList(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Auto migrate customer table
	if err := global.GVA_DB.AutoMigrate(&example.ExaCustomer{}); err != nil {
		t.Fatalf("Failed to migrate customer table: %v", err)
	}

	// Create multiple test customers for list
	testCustomers := []example.ExaCustomer{
		{
			CustomerName:       "Customer 1",
			CustomerPhoneData: "1234567890",
			SysUserID:         1,
			SysUserAuthorityID: 1,
		},
		{
			CustomerName:       "Customer 2",
			CustomerPhoneData: "0987654321",
			SysUserID:         2,
			SysUserAuthorityID: 1,
		},
	}
	for _, customer := range testCustomers {
		if err := global.GVA_DB.Create(&customer).Error; err != nil {
			t.Fatalf("Failed to create test customer: %v", err)
		}
	}

	tests := []struct {
		name                string
		sysUserAuthorityID uint
		info               request.PageInfo
		wantErr           bool
		checkResults      bool
	}{
		{
			name:                "Test Get Customer List Success",
			sysUserAuthorityID: 1,
			info:               request.PageInfo{Page: 1, PageSize: 10},
			wantErr:           false,
			checkResults:      true,
		},
		{
			name:                "Test Get Customer List Invalid Page",
			sysUserAuthorityID: 1,
			info:               request.PageInfo{Page: 0, PageSize: 10},
			wantErr:           true,
			checkResults:      false,
		},
		{
			name:                "Test Get Customer List Invalid PageSize",
			sysUserAuthorityID: 1,
			info:               request.PageInfo{Page: 1, PageSize: 0},
			wantErr:           true,
			checkResults:      false,
		},
		{
			name:                "Test Get Customer List Invalid Authority",
			sysUserAuthorityID: 999,
			info:               request.PageInfo{Page: 1, PageSize: 10},
			wantErr:           true,
			checkResults:      false,
		},
		{
			name:                "Test Get Customer List Large Page",
			sysUserAuthorityID: 1,
			info:               request.PageInfo{Page: 999999, PageSize: 10},
			wantErr:           false,
			checkResults:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CustomerService{}
			list, total, err := service.GetCustomerInfoList(tt.sysUserAuthorityID, tt.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCustomerInfoList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkResults {
				if list == nil {
					t.Error("GetCustomerInfoList() list should not be nil")
				}
				if total < 0 {
					t.Errorf("GetCustomerInfoList() total = %v, should not be negative", total)
				}

				// Type assertion and specific checks
				if customerList, ok := list.([]example.ExaCustomer); ok {
					// Verify the length doesn't exceed page size
					if len(customerList) > tt.info.PageSize {
						t.Errorf("GetCustomerInfoList() returned %v items, want <= %v", len(customerList), tt.info.PageSize)
					}
				} else {
					t.Error("GetCustomerInfoList() list is not of type []example.ExaCustomer")
				}
			}
		})
	}
}
