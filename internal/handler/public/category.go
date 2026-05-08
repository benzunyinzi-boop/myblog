package public

import (
	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/response"
	"github.com/yinyin/myblog/internal/service"
	"go.uber.org/zap"
)

// CategoryHandler 公开分类
type CategoryHandler struct {
	svc service.CategoryService
}

// NewCategoryHandler 构造
func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// List GET /api/v1/public/categories
func (h *CategoryHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		logger.C(c.Request.Context()).Error("public list categories", zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
		return
	}
	response.OK(c, gin.H{"items": list, "total": len(list)})
}
