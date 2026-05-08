package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/response"
	"github.com/yinyin/myblog/internal/pkg/uploader"
	"go.uber.org/zap"
)

// UploadHandler 管理端文件上传
type UploadHandler struct {
	uploader uploader.Uploader
}

// NewUploadHandler 构造
func NewUploadHandler(u uploader.Uploader) *UploadHandler {
	return &UploadHandler{uploader: u}
}

// Upload POST /api/v1/admin/uploads
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailInvalid(c, err)
		return
	}

	// 限制文件大小(10MB)
	const maxSize = 10 << 20
	if file.Size > maxSize {
		response.Fail(c, errcode.ErrInvalidParam, "文件大小不能超过 10MB")
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		logger.C(c.Request.Context()).Error("open upload file", zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
		return
	}
	defer src.Close()

	// 上传
	url, err := h.uploader.Upload(c.Request.Context(), file.Filename, src)
	if err != nil {
		logger.C(c.Request.Context()).Error("upload file", zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
		return
	}

	response.OK(c, gin.H{"url": url})
}
