package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"go.uber.org/zap"
)

// AccessLog 结构化访问日志(CLAUDE.md 要求:trace_id/path/method/status/latency)
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		logger.C(c.Request.Context()).Info("http_access",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("size", c.Writer.Size()),
		)
	}
}
