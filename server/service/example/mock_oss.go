package example

import (
	"mime/multipart"
)

type MockOss struct {
	shouldFail bool
}

func NewMockOss(shouldFail bool) *MockOss {
	return &MockOss{shouldFail: shouldFail}
}

func (m *MockOss) DeleteFile(key string) error {
	if m.shouldFail {
		return errors.New("mock: failed to delete file")
	}
	return nil
}

func (m *MockOss) UploadFile(file *multipart.FileHeader) (string, string, error) {
	if m.shouldFail {
		return "", "", errors.New("mock: failed to upload file")
	}
	return "http://example.com/test.txt", "test/test.txt", nil
}
