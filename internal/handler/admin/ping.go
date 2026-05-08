package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/middleware"
	"github.com/yinyin/myblog/internal/pkg/response"
)

// Ping GET /api/v1/admin/ping — 验证 JWT 中间件工作,返回当前登录用户。
func Ping(c *gin.Context) {
	uid, _ := middleware.UserIDFrom(c)
	username, _ := c.Get(middleware.CtxUsername)
	role, _ := c.Get(middleware.CtxRole)
	response.OK(c, gin.H{
		"message":  "pong",
		"user_id":  uid,
		"username": username,
		"role":     role,
	})
}
