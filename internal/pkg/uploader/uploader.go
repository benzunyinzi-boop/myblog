// Package uploader 文件上传抽象层。
// 本地实现存到 uploads/,预留 OSS 扩展点。
package uploader

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Uploader 文件上传接口
type Uploader interface {
	// Upload 上传文件,返回访问 URL
	Upload(ctx context.Context, filename string, content io.Reader) (string, error)
}

// LocalUploader 本地文件系统实现
type LocalUploader struct {
	baseDir string // 存储根目录,如 ./uploads
	baseURL string // 访问 URL 前缀,如 /uploads
}

// NewLocalUploader 构造本地上传器
func NewLocalUploader(baseDir, baseURL string) (*LocalUploader, error) {
	// 确保目录存在
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}
	return &LocalUploader{baseDir: baseDir, baseURL: baseURL}, nil
}

// Upload 上传到本地文件系统,按日期分目录存储
func (u *LocalUploader) Upload(ctx context.Context, filename string, content io.Reader) (string, error) {
	// 按日期分目录:uploads/2026/05/08/uuid-filename.ext
	now := time.Now()
	dateDir := filepath.Join(u.baseDir, now.Format("2006/01/02"))
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		return "", fmt.Errorf("create date dir: %w", err)
	}

	// 生成唯一文件名:uuid-原始文件名
	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]
	if len(base) > 50 {
		base = base[:50] // 限制长度
	}
	uniqueName := fmt.Sprintf("%s-%s%s", uuid.New().String()[:8], base, ext)
	fullPath := filepath.Join(dateDir, uniqueName)

	// 写入文件
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, content); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	// 返回访问 URL:/uploads/2026/01/02/uuid-filename.ext
	relPath := filepath.Join(now.Format("2006/01/02"), uniqueName)
	url := filepath.ToSlash(filepath.Join(u.baseURL, relPath))
	return url, nil
}
