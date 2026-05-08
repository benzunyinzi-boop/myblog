package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/dto"
	"github.com/yinyin/myblog/internal/middleware"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/response"
	"github.com/yinyin/myblog/internal/service"
	"go.uber.org/zap"
)

// ArticleHandler 管理端文章
type ArticleHandler struct {
	svc service.ArticleService
}

// NewArticleHandler 构造
func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// Create POST /api/v1/admin/articles
func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.ArticleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailInvalid(c, err)
		return
	}
	authorID, _ := middleware.UserIDFrom(c)
	if authorID == 0 {
		response.Fail(c, errcode.ErrUnauthorized, "")
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), authorID, req)
	if err != nil {
		h.writeError(c, err, "create")
		return
	}
	response.OK(c, resp)
}

// Update PUT /api/v1/admin/articles/:id
func (h *ArticleHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req dto.ArticleUpdateReq
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

// Delete DELETE /api/v1/admin/articles/:id
func (h *ArticleHandler) Delete(c *gin.Context) {
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

// Get GET /api/v1/admin/articles/:id
func (h *ArticleHandler) Get(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	resp, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err, "get")
		return
	}
	response.OK(c, resp)
}

// List GET /api/v1/admin/articles
func (h *ArticleHandler) List(c *gin.Context) {
	var req dto.ArticleListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailInvalid(c, err)
		return
	}
	resp, err := h.svc.List(c.Request.Context(), req)
	if err != nil {
		h.writeError(c, err, "list")
		return
	}
	response.OK(c, resp)
}

// Publish POST /api/v1/admin/articles/:id/publish
func (h *ArticleHandler) Publish(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	resp, err := h.svc.Publish(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err, "publish")
		return
	}
	response.OK(c, resp)
}

// Unpublish POST /api/v1/admin/articles/:id/unpublish
func (h *ArticleHandler) Unpublish(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	resp, err := h.svc.Unpublish(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err, "unpublish")
		return
	}
	response.OK(c, resp)
}

func (h *ArticleHandler) writeError(c *gin.Context, err error, op string) {
	switch {
	case errors.Is(err, service.ErrArticleNotFound):
		response.Fail(c, errcode.ErrArticleNotFound, "")
	case errors.Is(err, service.ErrArticleDuplicate):
		response.Fail(c, errcode.ErrArticleDuplicate, "")
	case errors.Is(err, service.ErrArticleConflict):
		response.Fail(c, errcode.ErrArticleConflict, "")
	case errors.Is(err, service.ErrCategoryNotFound):
		response.Fail(c, errcode.ErrCategoryNotFound, "")
	case errors.Is(err, service.ErrTagNotFound):
		response.Fail(c, errcode.ErrTagNotFound, "")
	default:
		logger.C(c.Request.Context()).Error("article op failed",
			zap.String("op", op), zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
	}
}
