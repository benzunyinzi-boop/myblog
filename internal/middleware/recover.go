package middleware

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/pkg/errcode"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/response"
	"go.uber.org/zap"
)

// Recover 捕获 panic,返回统一 500 响应
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.C(c.Request.Context()).Error("panic recovered",
					zap.Any("err", r),
					zap.String("stack", string(debug.Stack())),
				)
				response.Fail(c, errcode.ErrInternal, "")
				c.Abort()
			}
		}()
		c.Next()
	}
}
