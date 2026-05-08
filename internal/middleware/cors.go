package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 开发环境宽松 CORS。生产环境建议收紧 allow origins。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}
		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Credentials", "true")
		h.Set("Access-Control-Allow-Headers",
			"Content-Type,Authorization,X-Trace-Id,X-Requested-With")
		h.Set("Access-Control-Allow-Methods",
			"GET,POST,PUT,PATCH,DELETE,OPTIONS")
		h.Set("Access-Control-Expose-Headers", "X-Trace-Id")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
