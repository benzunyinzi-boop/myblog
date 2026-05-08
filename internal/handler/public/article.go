package public

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

// ArticleHandler 公开文章
type ArticleHandler struct {
	svc service.ArticleService
}

// NewArticleHandler 构造
func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// List GET /api/v1/public/articles
func (h *ArticleHandler) List(c *gin.Context) {
	var req dto.ArticleListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailInvalid(c, err)
		return
	}
	// 公开接口强制只返回已发布
	req.Status = "published"
	resp, err := h.svc.List(c.Request.Context(), req)
	if err != nil {
		logger.C(c.Request.Context()).Error("public list articles", zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
		return
	}
	response.OK(c, resp)
}

// GetBySlug GET /api/v1/public/articles/:slug
func (h *ArticleHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.FailInvalid(c, errors.New("slug required"))
		return
	}
	detail, err := h.svc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			response.Fail(c, errcode.ErrArticleNotFound, "")
			return
		}
		logger.C(c.Request.Context()).Error("public get article by slug", zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
		return
	}
	// 只返回已发布的文章
	if detail.Status != "published" {
		response.Fail(c, errcode.ErrArticleNotFound, "")
		return
	}
	// 异步增加浏览计数(不阻塞响应,失败也不影响)
	go func() {
		_ = h.svc.IncrView(c.Request.Context(), detail.ID)
	}()
	response.OK(c, detail)
}
