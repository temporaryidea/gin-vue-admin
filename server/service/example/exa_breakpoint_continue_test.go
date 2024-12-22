package example

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
)

func TestFindOrCreateFile(t *testing.T) {
	tests := []struct {
		name       string
		fileMd5    string
		fileName   string
		chunkTotal int
		wantErr    bool
		checkFile  bool // whether to check file properties
	}{
		{
			name:       "Test Create New File Success",
			fileMd5:    "test123md5",
			fileName:   "test.txt",
			chunkTotal: 3,
			wantErr:    false,
			checkFile:  true,
		},
		{
			name:       "Test Empty MD5",
			fileMd5:    "",
			fileName:   "test.txt",
			chunkTotal: 3,
			wantErr:    true,
			checkFile:  false,
		},
		{
			name:       "Test Empty Filename",
			fileMd5:    "test123md5",
			fileName:   "",
			chunkTotal: 3,
			wantErr:    true,
			checkFile:  false,
		},
		{
			name:       "Test Invalid Chunk Total",
			fileMd5:    "test123md5",
			fileName:   "test.txt",
			chunkTotal: 0,
			wantErr:    true,
			checkFile:  false,
		},
		{
			name:       "Test Duplicate MD5 Finished File",
			fileMd5:    "duplicate123",
			fileName:   "test2.txt",
			chunkTotal: 3,
			wantErr:    false,
			checkFile:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			file, err := service.FindOrCreateFile(tt.fileMd5, tt.fileName, tt.chunkTotal)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOrCreateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkFile {
				if file.FileMd5 != tt.fileMd5 {
					t.Errorf("FindOrCreateFile() file.FileMd5 = %v, want %v", file.FileMd5, tt.fileMd5)
				}
				if file.FileName != tt.fileName {
					t.Errorf("FindOrCreateFile() file.FileName = %v, want %v", file.FileName, tt.fileName)
				}
				if file.ChunkTotal != tt.chunkTotal {
					t.Errorf("FindOrCreateFile() file.ChunkTotal = %v, want %v", file.ChunkTotal, tt.chunkTotal)
				}
			}
		})
	}
}

func TestCreateFileChunk(t *testing.T) {
	tests := []struct {
		name            string
		id              uint
		fileChunkPath   string
		fileChunkNumber int
		wantErr         bool
	}{
		{
			name:            "Test Create File Chunk Success",
			id:              1,
			fileChunkPath:   "/tmp/chunk1",
			fileChunkNumber: 1,
			wantErr:         false,
		},
		{
			name:            "Test Invalid File ID",
			id:              0,
			fileChunkPath:   "/tmp/chunk1",
			fileChunkNumber: 1,
			wantErr:         true,
		},
		{
			name:            "Test Empty Chunk Path",
			id:              1,
			fileChunkPath:   "",
			fileChunkNumber: 1,
			wantErr:         true,
		},
		{
			name:            "Test Invalid Chunk Number",
			id:              1,
			fileChunkPath:   "/tmp/chunk1",
			fileChunkNumber: -1,
			wantErr:         true,
		},
		{
			name:            "Test Large Chunk Number",
			id:              1,
			fileChunkPath:   "/tmp/chunk1",
			fileChunkNumber: 1000000,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &FileUploadAndDownloadService{}
			err := service.CreateFileChunk(tt.id, tt.fileChunkPath, tt.fileChunkNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFileChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteFileChunk(t *testing.T) {
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
