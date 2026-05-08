package public

import (
	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/response"
	"github.com/yinyin/myblog/internal/service"
	"go.uber.org/zap"
)

// TagHandler 公开标签
type TagHandler struct{ svc service.TagService }

// NewTagHandler 构造
func NewTagHandler(svc service.TagService) *TagHandler { return &TagHandler{svc: svc} }

// List GET /api/v1/public/tags
func (h *TagHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		logger.C(c.Request.Context()).Error("public list tags", zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
		return
	}
	response.OK(c, gin.H{"items": list, "total": len(list)})
}
