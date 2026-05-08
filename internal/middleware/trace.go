// Package middleware 通用 HTTP 中间件。
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yinyin/myblog/internal/pkg/logger"
)

// TraceIDHeader 透传头
const TraceIDHeader = "X-Trace-Id"

// TraceID 生成或透传 trace_id,写入 ctx 与响应头
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := c.GetHeader(TraceIDHeader)
		if tid == "" {
			tid = uuid.NewString()
		}
		ctx := logger.WithTraceID(c.Request.Context(), tid)
		c.Request = c.Request.WithContext(ctx)
		c.Set("trace_id", tid)
		c.Writer.Header().Set(TraceIDHeader, tid)
		c.Next()
	}
}
