package example

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"go.uber.org/zap"
)

func TestUpload(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Auto migrate file upload table
	if err := global.GVA_DB.AutoMigrate(&example.ExaFileUploadAndDownload{}); err != nil {
		t.Fatalf("Failed to migrate file upload table: %v", err)
	}

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
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Auto migrate file upload table
	if err := global.GVA_DB.AutoMigrate(&example.ExaFileUploadAndDownload{}); err != nil {
		t.Fatalf("Failed to migrate file upload table: %v", err)
	}

	// Create test file record
	testFile := example.ExaFileUploadAndDownload{
		Name: "test.txt",
		Url:  "http://example.com/test.txt",
		Tag:  "txt",
		Key:  "test/test.txt",
	}
	if err := global.GVA_DB.Create(&testFile).Error; err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

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
				if file.GVA_MODEL.ID != tt.id {
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
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Set up test environment
	uploadPath := t.TempDir()
	testDir := filepath.Join(uploadPath, "test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Configure upload directory
	global.GVA_CONFIG.Local.StorePath = uploadPath
	global.GVA_CONFIG.Local.Path = "/"

	// Set up mock OSS and logger
	global.GVA_CONFIG.System.OssType = "mock"
	mockOss := NewMockOss(false)
	// Store the mock OSS instance to ensure it's reused
	upload.RegisterOssType("mock", func() upload.OSS {
		return mockOss
	})
	// Force initialization of OSS to use our mock
	if oss := upload.NewOss(); oss == nil {
		t.Fatal("Failed to initialize mock OSS")
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	global.GVA_LOG = logger

	defer func() {
		global.GVA_CONFIG.System.OssType = "local"
	}()

	// Auto migrate file upload table
	if migrateErr := global.GVA_DB.AutoMigrate(&example.ExaFileUploadAndDownload{}); migrateErr != nil {
		t.Fatalf("Failed to migrate file upload table: %v", migrateErr)
	}

	// Create test file
	testFilePath := filepath.Join(testDir, "test.txt")
	if err := os.WriteFile(testFilePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create test file record
	testFile := example.ExaFileUploadAndDownload{
		Name: "test.txt",
		Url:  filepath.Join(uploadPath, "test/test.txt"),
		Tag:  "txt",
		Key:  "test/test.txt", // Key must start with "test/" for OSS validation
	}
	// Create test file in database
	if err := global.GVA_DB.Create(&testFile).Error; err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	// Verify the file was created
	var createdFile example.ExaFileUploadAndDownload
	if err := global.GVA_DB.First(&createdFile, testFile.ID).Error; err != nil {
		t.Fatalf("Failed to verify test file creation: %v", err)
	}
	t.Logf("Created test file with ID: %d, Key: %s", createdFile.ID, createdFile.Key)

	tests := []struct {
		name    string
		file    example.ExaFileUploadAndDownload
		wantErr bool
	}{
		{
			name: "Test Delete Existing File",
			file: example.ExaFileUploadAndDownload{
				GVA_MODEL: global.GVA_MODEL{ID: testFile.ID}, // Use the actual ID from our created test file
				Name:      "test.txt",
				Key:       "test/test.txt",
			},
			wantErr: false,
		},
		{
			name: "Test Delete Non-existent File",
			file: example.ExaFileUploadAndDownload{
				GVA_MODEL: global.GVA_MODEL{ID: 999},
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
				GVA_MODEL: global.GVA_MODEL{ID: 1},
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
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Auto migrate file upload table
	if err := global.GVA_DB.AutoMigrate(&example.ExaFileUploadAndDownload{}); err != nil {
		t.Fatalf("Failed to migrate file upload table: %v", err)
	}

	// Create test file for editing
	testFile := example.ExaFileUploadAndDownload{
		Name: "test.txt",
		Url:  "http://example.com/test.txt",
		Tag:  "txt",
		Key:  "test/test.txt",
	}
	if err := global.GVA_DB.Create(&testFile).Error; err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		file    example.ExaFileUploadAndDownload
		wantErr bool
	}{
		{
			name: "Test Edit Existing File Name",
			file: example.ExaFileUploadAndDownload{
				GVA_MODEL: global.GVA_MODEL{ID: 1},
				Name: "updated.txt",
			},
			wantErr: false,
		},
		{
			name: "Test Edit Non-existent File Name",
			file: example.ExaFileUploadAndDownload{
				GVA_MODEL: global.GVA_MODEL{ID: 999},
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
				GVA_MODEL: global.GVA_MODEL{ID: 1},
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
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Ensure database connection
	if global.GVA_DB == nil {
		t.Fatal("Database connection is nil")
	}

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
			wantErr:      false,
			checkResults: true,
		},
		{
			name: "Test Get File List Invalid PageSize",
			info: request.PageInfo{
				Page:     1,
				PageSize: 0,
			},
			wantErr:      false,
			checkResults: true,
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
func createTestFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = part.Write(content)
	if err != nil {
		t.Fatalf("Failed to write content: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(1024)
	if err != nil {
		t.Fatalf("Failed to read form: %v", err)
	}

	header := form.File["file"][0]
	if header == nil {
		t.Fatal("Failed to create file header")
	}
	return header
}

func TestUploadFile(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Set up mock OSS and ensure upload directory exists
	global.GVA_CONFIG.System.OssType = "mock"
	upload.RegisterOssType("mock", func() upload.OSS {
		return NewMockOss(false)
	})

	// Configure upload directory
	uploadPath := t.TempDir()
	uploadDir := filepath.Join(uploadPath, "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		t.Fatalf("Failed to create upload directory: %v", err)
	}
	global.GVA_CONFIG.Local.StorePath = uploadDir
	global.GVA_CONFIG.Local.Path = "/uploads"

	// Create test files with different content types
	textContent := []byte("test content")
	textHeader := createTestFileHeader(t, "test.txt", textContent)
	imageContent := []byte{0xFF, 0xD8, 0xFF} // Simple JPEG header
	imageHeader := createTestFileHeader(t, "test.jpg", imageContent)

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
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Auto migrate file upload table
	if err := global.GVA_DB.AutoMigrate(&example.ExaFileUploadAndDownload{}); err != nil {
		t.Fatalf("Failed to migrate file upload table: %v", err)
	}

	tests := []struct {
		name    string
		urls    []string
		wantErr bool
	}{
		{
			name: "Test Import Multiple Valid URLs",
			urls: []string{
				"http://example.com/test1.txt",
				"http://example.com/test2.jpg",
			},
			wantErr: false,
		},
		{
			name: "Test Import Single Valid URL",
			urls: []string{
				"http://example.com/test.pdf",
			},
			wantErr: false,
		},
		{
			name:    "Test Import Nil URLs",
			urls:    nil,
			wantErr: true,
		},
		{
			name:    "Test Import Empty URLs",
			urls:    []string{},
			wantErr: false,
		},
		{
			name: "Test Import Invalid URLs",
			urls: []string{
				"",
				"http://example.com/",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear database before each test
			if err := global.GVA_DB.Where("1 = 1").Delete(&example.ExaFileUploadAndDownload{}).Error; err != nil {
				t.Fatalf("Failed to clear database: %v", err)
			}

			service := &FileUploadAndDownloadService{}
			err := service.ImportURL(tt.urls)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportURL() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.urls != nil {
				var files []example.ExaFileUploadAndDownload
				if err := global.GVA_DB.Find(&files).Error; err != nil {
					t.Errorf("Failed to retrieve imported files: %v", err)
				}

				// Check that valid URLs were imported
				for _, url := range tt.urls {
					if url == "" || !strings.Contains(url, "/") {
						continue
					}
					s := strings.Split(url, "/")
					if s[len(s)-1] == "" {
						continue
					}
					found := false
					for _, file := range files {
						if file.Url == url {
							found = true
							if file.Name != s[len(s)-1] {
								t.Errorf("ImportURL() file name = %v, want %v", file.Name, s[len(s)-1])
							}
							if file.Tag != "url" {
								t.Errorf("ImportURL() file tag = %v, want url", file.Tag)
							}
							if !strings.HasPrefix(file.Key, "test/") {
								t.Errorf("ImportURL() file key = %v, should start with test/", file.Key)
							}
							break
						}
					}
					if !found {
						t.Errorf("ImportURL() failed to import URL: %v", url)
					}
				}
			}
		})
	}
}
