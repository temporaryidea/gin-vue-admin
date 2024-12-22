package example

import (
	"reflect"
	"testing"
)

func TestServiceGroup(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "Test ServiceGroup Structure",
			test: func(t *testing.T) {
				// Test type structure
				group := ServiceGroup{}
				groupType := reflect.TypeOf(group)

				// Verify number of fields
				if groupType.NumField() != 2 {
					t.Errorf("ServiceGroup has %v fields, want 2", groupType.NumField())
				}

				// Verify CustomerService field
				if _, exists := groupType.FieldByName("CustomerService"); !exists {
					t.Error("ServiceGroup missing CustomerService field")
				}

				// Verify FileUploadAndDownloadService field
				if _, exists := groupType.FieldByName("FileUploadAndDownloadService"); !exists {
					t.Error("ServiceGroup missing FileUploadAndDownloadService field")
				}
			},
		},
		{
			name: "Test Service Composition",
			test: func(t *testing.T) {
				// Test ability to compose services
				group := ServiceGroup{
					CustomerService:              CustomerService{},
					FileUploadAndDownloadService: FileUploadAndDownloadService{},
				}

				// Verify service types
				if reflect.TypeOf(group.CustomerService) != reflect.TypeOf(CustomerService{}) {
					t.Error("CustomerService field has incorrect type")
				}

				if reflect.TypeOf(group.FileUploadAndDownloadService) != reflect.TypeOf(FileUploadAndDownloadService{}) {
					t.Error("FileUploadAndDownloadService field has incorrect type")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
