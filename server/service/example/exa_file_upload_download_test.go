package example

import (
	"bytes"
	"errors"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
)

func TestUpload(t *testing.T) {
	tests := []struct {
		name    string
		file    example.ExaFileUploadAndDownload
		wantErr bool
	}{
		{
			name: "Test Valid Upload",
			file: example.ExaFileUploadAndDownload{
				Name: "test.txt",
				Url:  "http://example.com/test.txt",
				Tag:  "txt",
				Key:  "test/test.txt",
			},
			wantErr: false,
		},
		{
			name: "Test Upload Without Name",
			file: example.ExaFileUploadAndDownload{
				Url: "http://example.com/test.txt",
				Tag: "txt",
				Key: "test/test.txt",
			},
			wantErr: true,
		},
		{
			name: "Test Upload Without URL",
			file: example.ExaFileUploadAndDownload{
				Name: "test.txt",
				Tag:  "txt",
				Key:  "test/test.txt",
			},
			wantErr: true,
		},
		{
			name: "Test Upload Without Key",
			file: example.ExaFileUploadAndDownload{
				Name: "test.txt",
				Url:  "http://example.com/test.txt",
				Tag:  "txt",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			err := service.Upload(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindFile(t *testing.T) {
	tests := []struct {
		name       string
		id         uint
		wantErr    bool
		checkProps bool
	}{
		{
			name:       "Test Find Existing File",
			id:         1,
			wantErr:    false,
			checkProps: true,
		},
		{
			name:       "Test Find Non-existent File",
			id:         999,
			wantErr:    true,
			checkProps: false,
		},
		{
			name:       "Test Find File With Zero ID",
			id:         0,
			wantErr:    true,
			checkProps: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			file, err := service.FindFile(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkProps {
				if file.ID != tt.id {
					t.Errorf("FindFile() file.ID = %v, want %v", file.ID, tt.id)
				}
				if file.Name == "" {
					t.Error("FindFile() file.Name should not be empty")
				}
				if file.Url == "" {
					t.Error("FindFile() file.Url should not be empty")
				}
				if file.Key == "" {
					t.Error("FindFile() file.Key should not be empty")
				}
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		name    string
		file    example.ExaFileUploadAndDownload
		wantErr bool
	}{
		{
			name: "Test Delete Existing File",
			file: example.ExaFileUploadAndDownload{
				ID:   1,
				Name: "test.txt",
				Key:  "test/test.txt",
			},
			wantErr: false,
		},
		{
			name: "Test Delete Non-existent File",
			file: example.ExaFileUploadAndDownload{
				ID:   999,
				Name: "nonexistent.txt",
				Key:  "test/nonexistent.txt",
			},
			wantErr: true,
		},
		{
			name: "Test Delete File Without ID",
			file: example.ExaFileUploadAndDownload{
				Name: "test.txt",
				Key:  "test/test.txt",
			},
			wantErr: true,
		},
		{
			name: "Test Delete File Without Key",
			file: example.ExaFileUploadAndDownload{
				ID:   1,
				Name: "test.txt",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			err := service.DeleteFile(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEditFileName(t *testing.T) {
	tests := []struct {
		name    string
		file    example.ExaFileUploadAndDownload
		wantErr bool
	}{
		{
			name: "Test Edit Existing File Name",
			file: example.ExaFileUploadAndDownload{
				ID:   1,
				Name: "updated.txt",
			},
			wantErr: false,
		},
		{
			name: "Test Edit Non-existent File Name",
			file: example.ExaFileUploadAndDownload{
				ID:   999,
				Name: "nonexistent.txt",
			},
			wantErr: true,
		},
		{
			name: "Test Edit File Without ID",
			file: example.ExaFileUploadAndDownload{
				Name: "test.txt",
			},
			wantErr: true,
		},
		{
			name: "Test Edit File With Empty Name",
			file: example.ExaFileUploadAndDownload{
				ID:   1,
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			err := service.EditFileName(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("EditFileName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFileRecordInfoList(t *testing.T) {
	tests := []struct {
		name         string
		info         request.PageInfo
		wantErr      bool
		checkResults bool
	}{
		{
			name: "Test Get File List Success",
			info: request.PageInfo{
				Page:     1,
				PageSize: 10,
				Keyword:  "",
			},
			wantErr:      false,
			checkResults: true,
		},
		{
			name: "Test Get File List With Keyword",
			info: request.PageInfo{
				Page:     1,
				PageSize: 10,
				Keyword:  "test",
			},
			wantErr:      false,
			checkResults: true,
		},
		{
			name: "Test Get File List Invalid Page",
			info: request.PageInfo{
				Page:     0,
				PageSize: 10,
			},
			wantErr:      true,
			checkResults: false,
		},
		{
			name: "Test Get File List Invalid PageSize",
			info: request.PageInfo{
				Page:     1,
				PageSize: 0,
			},
			wantErr:      true,
			checkResults: false,
		},
		{
			name: "Test Get File List Large Page",
			info: request.PageInfo{
				Page:     999999,
				PageSize: 10,
			},
			wantErr:      false,
			checkResults: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			list, total, err := service.GetFileRecordInfoList(tt.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileRecordInfoList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkResults {
				if list == nil {
					t.Error("GetFileRecordInfoList() list should not be nil")
				}
				if total < 0 {
					t.Errorf("GetFileRecordInfoList() total = %v, should not be negative", total)
				}

				// Type assertion and specific checks
				if fileList, ok := list.([]example.ExaFileUploadAndDownload); ok {
					// Verify the length doesn't exceed page size
					if len(fileList) > tt.info.PageSize {
						t.Errorf("GetFileRecordInfoList() returned %v items, want <= %v", len(fileList), tt.info.PageSize)
					}

					// Check keyword filter if provided
					if tt.info.Keyword != "" {
						for _, file := range fileList {
							if !strings.Contains(file.Name, tt.info.Keyword) {
								t.Errorf("GetFileRecordInfoList() file name %v does not contain keyword %v", file.Name, tt.info.Keyword)
							}
						}
					}
				} else {
					t.Error("GetFileRecordInfoList() list is not of type []example.ExaFileUploadAndDownload")
				}
			}
		})
	}
}

// createTestFileHeader creates a multipart.FileHeader for testing
func createTestFileHeader(filename string, content []byte) (*multipart.FileHeader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(content)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(1024)
	if err != nil {
		return nil, err
	}

	return form.File["file"][0], nil
}

func TestUploadFile(t *testing.T) {
	// Create test files with different content types
	textContent := []byte("test content")
	textHeader, err := createTestFileHeader("test.txt", textContent)
	if err != nil {
		t.Fatalf("Failed to create text file header: %v", err)
	}

	imageContent := []byte{0xFF, 0xD8, 0xFF} // Simple JPEG header
	imageHeader, err := createTestFileHeader("test.jpg", imageContent)
	if err != nil {
		t.Fatalf("Failed to create image file header: %v", err)
	}

	tests := []struct {
		name       string
		header     *multipart.FileHeader
		noSave     string
		wantErr    bool
		checkProps bool
	}{
		{
			name:       "Test Upload Text File With Save",
			header:     textHeader,
			noSave:     "0",
			wantErr:    false,
			checkProps: true,
		},
		{
			name:       "Test Upload Image File With Save",
			header:     imageHeader,
			noSave:     "0",
			wantErr:    false,
			checkProps: true,
		},
		{
			name:       "Test Upload Without Save",
			header:     textHeader,
			noSave:     "1",
			wantErr:    false,
			checkProps: true,
		},
		{
			name:       "Test Upload Nil Header",
			header:     nil,
			noSave:     "0",
			wantErr:    true,
			checkProps: false,
		},
		{
			name:       "Test Upload Invalid NoSave Value",
			header:     textHeader,
			noSave:     "invalid",
			wantErr:    true,
			checkProps: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			file, err := service.UploadFile(tt.header, tt.noSave)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkProps && !tt.wantErr {
				if file.Name != tt.header.Filename {
					t.Errorf("UploadFile() file.Name = %v, want %v", file.Name, tt.header.Filename)
				}
				if file.Tag == "" {
					t.Error("UploadFile() file.Tag should not be empty")
				}
				if file.Key == "" {
					t.Error("UploadFile() file.Key should not be empty")
				}
				if file.Url == "" {
					t.Error("UploadFile() file.Url should not be empty")
				}

				// Check file extension handling
				ext := strings.Split(tt.header.Filename, ".")[1]
				if file.Tag != ext {
					t.Errorf("UploadFile() file.Tag = %v, want %v", file.Tag, ext)
				}
			}
		})
	}
}

func TestImportURL(t *testing.T) {
	tests := []struct {
		name    string
		file    *[]example.ExaFileUploadAndDownload
		wantErr bool
	}{
		{
			name: "Test Import Multiple Valid URLs",
			file: &[]example.ExaFileUploadAndDownload{
				{
					Name: "test1.txt",
					Url:  "http://example.com/test1.txt",
					Tag:  "txt",
					Key:  "test/test1.txt",
				},
				{
					Name: "test2.jpg",
					Url:  "http://example.com/test2.jpg",
					Tag:  "jpg",
					Key:  "test/test2.jpg",
				},
			},
			wantErr: false,
		},
		{
			name: "Test Import Single Valid URL",
			file: &[]example.ExaFileUploadAndDownload{
				{
					Name: "test.pdf",
					Url:  "http://example.com/test.pdf",
					Tag:  "pdf",
					Key:  "test/test.pdf",
				},
			},
			wantErr: false,
		},
		{
			name:    "Test Import Nil URLs",
			file:    nil,
			wantErr: true,
		},
		{
			name:    "Test Import Empty URLs",
			file:    &[]example.ExaFileUploadAndDownload{},
			wantErr: true,
		},
		{
			name: "Test Import Invalid URLs",
			file: &[]example.ExaFileUploadAndDownload{
				{
					Name: "",
					Url:  "",
					Tag:  "",
					Key:  "",
				},
			},
			wantErr: true,
		},
		{
			name: "Test Import Mixed Valid and Invalid URLs",
			file: &[]example.ExaFileUploadAndDownload{
				{
					Name: "valid.txt",
					Url:  "http://example.com/valid.txt",
					Tag:  "txt",
					Key:  "test/valid.txt",
				},
				{
					Name: "",
					Url:  "",
					Tag:  "",
					Key:  "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			err := service.ImportURL(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportURL() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Additional checks for successful imports
			if !tt.wantErr && tt.file != nil {
				for _, f := range *tt.file {
					if f.Name == "" {
						t.Error("ImportURL() file.Name should not be empty")
					}
					if f.Url == "" {
						t.Error("ImportURL() file.Url should not be empty")
					}
					if f.Tag == "" {
						t.Error("ImportURL() file.Tag should not be empty")
					}
					if f.Key == "" {
						t.Error("ImportURL() file.Key should not be empty")
					}
				}
			}
		})
	}
}
