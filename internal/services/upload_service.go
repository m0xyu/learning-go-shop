package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/m0xyu/learning-go-shop/internal/interfaces"
)

type UploadService struct {
	provider interfaces.UploadProvider
}

func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{provider: provider}
}

func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	// ファイルの拡張子をチェックする
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExt(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	newfileName := uuid.New().String()

	path := fmt.Sprintf("products/%d/%s%s", productID, newfileName, ext)

	return s.provider.UploadFile(file, path)
}

func isValidImageExt(ext string) bool {
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, v := range validExts {
		if ext == v {
			return true
		}
	}
	return false
}
