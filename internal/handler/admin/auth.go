// Package admin 管理端(需鉴权)handler。
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

// AuthHandler 管理端鉴权 handler
type AuthHandler struct {
	auth service.AuthService
}

// NewAuthHandler 构造
func NewAuthHandler(auth service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Login POST /api/v1/admin/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailInvalid(c, err)
		return
	}

	resp, err := h.auth.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPassword),
			errors.Is(err, service.ErrUserNotFound):
			response.Fail(c, errcode.ErrInvalidPassword, "")
		case errors.Is(err, service.ErrUserDisabled):
			response.Fail(c, errcode.ErrUserDisabled, "")
		default:
			logger.C(c.Request.Context()).Error("login failed", zap.Error(err), zap.String("username", req.Username))
			response.Fail(c, errcode.ErrInternal, "")
		}
		return
	}
	response.OK(c, resp)
}
