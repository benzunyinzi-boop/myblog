package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/dto"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/response"
	"github.com/yinyin/myblog/internal/service"
	"go.uber.org/zap"
)

// ProfileHandler 管理端个人资料
type ProfileHandler struct {
	svc service.ProfileService
}

// NewProfileHandler 构造
func NewProfileHandler(svc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// Update PUT /api/v1/admin/profile
func (h *ProfileHandler) Update(c *gin.Context) {
	var req dto.ProfileUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailInvalid(c, err)
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), req)
	if err != nil {
		logger.C(c.Request.Context()).Error("profile update failed", zap.Error(err))
		response.Fail(c, errcode.ErrInternal, "")
		return
	}
	response.OK(c, resp)
}
