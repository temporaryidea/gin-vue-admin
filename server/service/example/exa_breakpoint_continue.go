package example

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"gorm.io/gorm"
)

type FileUploadAndDownloadService struct{}

var FileUploadAndDownloadServiceApp = new(FileUploadAndDownloadService)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: FindOrCreateFile
//@description: 上传文件时检测当前文件属性，如果没有文件则创建，有则返回文件的当前切片
//@param: fileMd5 string, fileName string, chunkTotal int
//@return: file model.ExaFile, err error

func (e *FileUploadAndDownloadService) FindOrCreateFile(fileMd5 string, fileName string, chunkTotal int) (file example.ExaFile, err error) {
	if fileMd5 == "" {
		return file, errors.New("file md5 cannot be empty")
	}
	if fileName == "" {
		return file, errors.New("file name cannot be empty")
	}
	if chunkTotal <= 0 {
		return file, errors.New("chunk total must be positive")
	}

	var cfile example.ExaFile
	cfile.FileMd5 = fileMd5
	cfile.FileName = fileName
	cfile.ChunkTotal = chunkTotal

	if errors.Is(global.GVA_DB.Where("file_md5 = ? AND is_finish = ?", fileMd5, true).First(&file).Error, gorm.ErrRecordNotFound) {
		err = global.GVA_DB.Where("file_md5 = ? AND file_name = ?", fileMd5, fileName).Preload("ExaFileChunk").FirstOrCreate(&file, cfile).Error
		return file, err
	}
	cfile.IsFinish = true
	cfile.FilePath = file.FilePath
	err = global.GVA_DB.Create(&cfile).Error
	return cfile, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateFileChunk
//@description: 创建文件切片记录
//@param: id uint, fileChunkPath string, fileChunkNumber int
//@return: error

func (e *FileUploadAndDownloadService) CreateFileChunk(id uint, fileChunkPath string, fileChunkNumber int) error {
	if id == 0 {
		return errors.New("invalid file ID")
	}
	if fileChunkPath == "" {
		return errors.New("chunk path cannot be empty")
	}
	if fileChunkNumber <= 0 {
		return errors.New("chunk number must be positive")
	}
	if fileChunkNumber > 100000 { // reasonable limit for chunk numbers
		return errors.New("chunk number too large")
	}

	var chunk example.ExaFileChunk
	chunk.FileChunkPath = fileChunkPath
	chunk.ExaFileID = id
	chunk.FileChunkNumber = fileChunkNumber
	err := global.GVA_DB.Create(&chunk).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFileChunk
//@description: 删除文件切片记录
//@param: fileMd5 string, fileName string, filePath string
//@return: error

func (e *FileUploadAndDownloadService) DeleteFileChunk(fileMd5 string, filePath string) error {
	if fileMd5 == "" {
		return errors.New("file md5 cannot be empty")
	}
	if filePath == "" {
		return errors.New("file path cannot be empty")
	}

	var chunks []example.ExaFileChunk
	var file example.ExaFile
	err := global.GVA_DB.Where("file_md5 = ?", fileMd5).First(&file).Error
	if err != nil {
		return err
	}
	err = global.GVA_DB.Model(&file).Updates(map[string]interface{}{
		"is_finish":  true,
		"file_path": filePath,
	}).Error
	if err != nil {
		return err
	}
	err = global.GVA_DB.Where("exa_file_id = ?", file.ID).Delete(&chunks).Unscoped().Error
	return err
}
