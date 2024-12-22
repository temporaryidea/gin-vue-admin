package example

import (
	"errors"
	"mime/multipart"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
)

type MockOss struct {
	shouldFail bool
}

// Verify interface implementation
var _ upload.OSS = (*MockOss)(nil)

func NewMockOss(shouldFail bool) *MockOss {
	return &MockOss{shouldFail: shouldFail}
}

func (m *MockOss) DeleteFile(key string) error {
	if m.shouldFail {
		return errors.New("mock: failed to delete file")
	}
	// For test purposes, we'll accept any key that starts with "test/"
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if !strings.HasPrefix(key, "test/") {
		return errors.New("非法的key")
	}
	// Mock successful deletion for valid keys
	return nil
}

func (m *MockOss) UploadFile(file *multipart.FileHeader) (string, string, error) {
	if m.shouldFail {
		return "", "", errors.New("mock: failed to upload file")
	}
	if file == nil {
		return "", "", errors.New("file is nil")
	}
	// Generate mock URL and key based on filename
	filename := file.Filename
	key := "test/" + filename
	url := "http://example.com/" + filename
	return url, key, nil
}
