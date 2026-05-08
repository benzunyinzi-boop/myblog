package admin

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/dto"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/response"
	"github.com/yinyin/myblog/internal/service"
	"go.uber.org/zap"
)

// CategoryHandler 管理端分类
type CategoryHandler struct {
	svc service.CategoryService
}

// NewCategoryHandler 构造
func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// Create POST /api/v1/admin/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryCreateReq
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

// Update PUT /api/v1/admin/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req dto.CategoryUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailInvalid(c, err)
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		h.writeError(c, err, "update")
		return
	}
	response.OK(c, resp)
}

// Delete DELETE /api/v1/admin/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
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

// List GET /api/v1/admin/categories
func (h *CategoryHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		h.writeError(c, err, "list")
		return
	}
	response.OK(c, gin.H{"items": list, "total": len(list)})
}

func (h *CategoryHandler) writeError(c *gin.Context, err error, op string) {
	switch {
	case errors.Is(err, service.ErrCategoryNotFound):
		response.Fail(c, errcode.ErrCategoryNotFound, "")
	case errors.Is(err, service.ErrCategoryDuplicate):
		response.Fail(c, errcode.ErrCategoryDuplicate, "")
	default:
		logger.C(c.Request.Context()).Error("category op failed",
			zap.String("op", op), zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
	}
}

// parseIDParam 解析 /:name 参数为 uint64,失败时写响应并返回 ok=false
func parseIDParam(c *gin.Context, name string) (uint64, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		response.FailInvalid(c, errors.New("invalid "+name))
		return 0, false
	}
	return id, true
}
