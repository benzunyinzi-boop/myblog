package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	myjwt "github.com/yinyin/myblog/internal/pkg/jwt"
	"github.com/yinyin/myblog/internal/pkg/response"
)

// Gin context 里的鉴权 key
const (
	CtxUserID   = "auth.user_id"
	CtxUsername = "auth.username"
	CtxRole     = "auth.role"
)

// JWTAuth 强制鉴权中间件,校验 Authorization: Bearer <token>
func JWTAuth(signer *myjwt.Signer) gin.HandlerFunc {
	return func(c *gin.Context) {
		tok := extractBearer(c.GetHeader("Authorization"))
		if tok == "" {
			response.Fail(c, errcode.ErrUnauthorized, "")
			c.Abort()
			return
		}
		claims, err := signer.Verify(tok, myjwt.TypeAccess)
		if err != nil {
			switch {
			case errors.Is(err, myjwt.ErrExpiredToken):
				response.Fail(c, errcode.ErrTokenExpired, "")
			default:
				response.Fail(c, errcode.ErrTokenInvalid, "")
			}
			c.Abort()
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func extractBearer(h string) string {
	const prefix = "Bearer "
	if len(h) <= len(prefix) {
		return ""
	}
	if !strings.EqualFold(h[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}

// UserIDFrom 从 gin.Context 提取当前登录用户 id。
func UserIDFrom(c *gin.Context) (uint64, bool) {
	v, ok := c.Get(CtxUserID)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}
