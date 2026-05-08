package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/dto"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/response"
	"github.com/yinyin/myblog/internal/service"
	"go.uber.org/zap"
)

// TagHandler 管理端标签
type TagHandler struct {
	svc service.TagService
}

// NewTagHandler 构造
func NewTagHandler(svc service.TagService) *TagHandler { return &TagHandler{svc: svc} }

// Create POST /api/v1/admin/tags
func (h *TagHandler) Create(c *gin.Context) {
	var req dto.TagCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailInvalid(c, err)
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		h.writeError(c, err, "create")
		return
	}
	response.OK(c, resp)
}

// Delete DELETE /api/v1/admin/tags/:id
func (h *TagHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.writeError(c, err, "delete")
		return
	}
	response.OK(c, nil)
}

// List GET /api/v1/admin/tags
func (h *TagHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		h.writeError(c, err, "list")
		return
	}
	response.OK(c, gin.H{"items": list, "total": len(list)})
}

func (h *TagHandler) writeError(c *gin.Context, err error, op string) {
	switch {
	case errors.Is(err, service.ErrTagNotFound):
		response.Fail(c, errcode.ErrTagNotFound, "")
	case errors.Is(err, service.ErrTagDuplicate):
		response.Fail(c, errcode.ErrTagDuplicate, "")
	default:
		logger.C(c.Request.Context()).Error("tag op failed",
			zap.String("op", op), zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
	}
}
