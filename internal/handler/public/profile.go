package public

import (
	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/response"
	"github.com/yinyin/myblog/internal/service"
	"go.uber.org/zap"
)

// ProfileHandler 公开个人资料
type ProfileHandler struct {
	svc service.ProfileService
}

// NewProfileHandler 构造
func NewProfileHandler(svc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// Get GET /api/v1/public/profile
func (h *ProfileHandler) Get(c *gin.Context) {
	resp, err := h.svc.Get(c.Request.Context())
	if err != nil {
		logger.C(c.Request.Context()).Error("public get profile", zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
		return
	}
	response.OK(c, resp)
}
