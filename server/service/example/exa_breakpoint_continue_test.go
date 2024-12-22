package example

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
)

func TestFindOrCreateFile(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()
	tests := []struct {
		name       string
		file       example.ExaFile
		chunkTotal int
		wantErr    bool
		checkFile  bool // whether to check file properties
	}{
		{
			name:       "Test Create New File Success",
			file: example.ExaFile{
				FileMd5:  "test123md5",
				FileName: "test.txt",
			},
			chunkTotal: 3,
			wantErr:    false,
			checkFile:  true,
		},
		{
			name:       "Test Empty MD5",
			file: example.ExaFile{
				FileMd5:  "",
				FileName: "test.txt",
			},
			chunkTotal: 3,
			wantErr:    true,
			checkFile:  false,
		},
		{
			name:       "Test Empty Filename",
			file: example.ExaFile{
				FileMd5:  "test123md5",
				FileName: "",
			},
			chunkTotal: 3,
			wantErr:    true,
			checkFile:  false,
		},
		{
			name:       "Test Invalid Chunk Total",
			file: example.ExaFile{
				FileMd5:  "test123md5",
				FileName: "test.txt",
			},
			chunkTotal: 0,
			wantErr:    true,
			checkFile:  false,
		},
		{
			name:       "Test Duplicate MD5 Finished File",
			file: example.ExaFile{
				FileMd5:  "duplicate123",
				FileName: "test2.txt",
			},
			chunkTotal: 3,
			wantErr:    false,
			checkFile:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			file, err := service.FindOrCreateFile(tt.file.FileMd5, tt.file.FileName, tt.chunkTotal)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOrCreateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkFile {
				if file.FileMd5 != tt.file.FileMd5 {
					t.Errorf("FindOrCreateFile() file.FileMd5 = %v, want %v", file.FileMd5, tt.file.FileMd5)
				}
				if file.FileName != tt.file.FileName {
					t.Errorf("FindOrCreateFile() file.FileName = %v, want %v", file.FileName, tt.file.FileName)
				}
				if file.ChunkTotal != tt.chunkTotal {
					t.Errorf("FindOrCreateFile() file.ChunkTotal = %v, want %v", file.ChunkTotal, tt.chunkTotal)
				}
			}
		})
	}
}

func TestCreateFileChunk(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()
	tests := []struct {
		name            string
		file            example.ExaFile
		fileChunkPath   string
		fileChunkNumber int
		wantErr         bool
	}{
		{
			name:            "Test Create File Chunk Success",
			file:            example.ExaFile{GVA_MODEL: global.GVA_MODEL{ID: 1}},
			fileChunkPath:   "/tmp/chunk1",
			fileChunkNumber: 1,
			wantErr:         false,
		},
		{
			name:            "Test Invalid File ID",
			file:            example.ExaFile{GVA_MODEL: global.GVA_MODEL{ID: 0}},
			fileChunkPath:   "/tmp/chunk1",
			fileChunkNumber: 1,
			wantErr:         true,
		},
		{
			name:            "Test Empty Chunk Path",
			file:            example.ExaFile{GVA_MODEL: global.GVA_MODEL{ID: 1}},
			fileChunkPath:   "",
			fileChunkNumber: 1,
			wantErr:         true,
		},
		{
			name:            "Test Invalid Chunk Number",
			file:            example.ExaFile{GVA_MODEL: global.GVA_MODEL{ID: 1}},
			fileChunkPath:   "/tmp/chunk1",
			fileChunkNumber: -1,
			wantErr:         true,
		},
		{
			name:            "Test Large Chunk Number",
			file:            example.ExaFile{GVA_MODEL: global.GVA_MODEL{ID: 1}},
			fileChunkPath:   "/tmp/chunk1",
			fileChunkNumber: 1000000,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			err := service.CreateFileChunk(tt.file.ID, tt.fileChunkPath, tt.fileChunkNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFileChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteFileChunk(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a test file for deletion test
	testFile := example.ExaFile{
		FileMd5:    "test123md5",
		FileName:   "test.txt",
		ChunkTotal: 5,
	}
	if err := global.GVA_DB.Create(&testFile).Error; err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create test chunks
	testChunk := example.ExaFileChunk{
		ExaFileID:       testFile.ID,
		FileChunkNumber: 1,
		FileChunkPath:   "/tmp/chunk1",
	}
	if err := global.GVA_DB.Create(&testChunk).Error; err != nil {
		t.Fatalf("Failed to create test chunk: %v", err)
	}

	tests := []struct {
		name     string
		fileMd5  string
		filePath string
		wantErr  bool
	}{
		{
			name:     "Test Delete File Chunk Success",
			fileMd5:  "test123md5",
			filePath: "/tmp/test.txt",
			wantErr:  false,
		},
		{
			name:     "Test Empty MD5",
			fileMd5:  "",
			filePath: "/tmp/test.txt",
			wantErr:  true,
		},
		{
			name:     "Test Empty File Path",
			fileMd5:  "test123md5",
			filePath: "",
			wantErr:  true,
		},
		{
			name:     "Test Non-existent MD5",
			fileMd5:  "nonexistent123",
			filePath: "/tmp/test.txt",
			wantErr:  true,
		},
		{
			name:     "Test Invalid File Path",
			fileMd5:  "test123md5",
			filePath: "/invalid/path/test.txt",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			err := service.DeleteFileChunk(tt.fileMd5, tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFileChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
