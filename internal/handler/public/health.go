// Package public 公开接口(无需鉴权)handler。
package public

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/pkg/response"
)

// HealthHandler 探活接口
type HealthHandler struct {
	service string
	version string
}

// NewHealthHandler 构造
func NewHealthHandler(service, version string) *HealthHandler {
	return &HealthHandler{service: service, version: version}
}

// Check GET /health
// 返回统一响应:{code:0, message:"ok", data:{status,service,version,ts}}
func (h *HealthHandler) Check(c *gin.Context) {
	response.OK(c, gin.H{
		"status":  "up",
		"service": h.service,
		"version": h.version,
		"ts":      time.Now().Unix(),
	})
}
