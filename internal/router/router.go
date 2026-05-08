// Package router 统一路由注册。
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/handler/admin"
	"github.com/yinyin/myblog/internal/handler/public"
	"github.com/yinyin/myblog/internal/middleware"
	myjwt "github.com/yinyin/myblog/internal/pkg/jwt"
	"github.com/yinyin/myblog/internal/service"
)

// Deps 路由注册依赖
type Deps struct {
	ServiceName string
	Version     string
	Signer      *myjwt.Signer       // 可为空,仅在启用鉴权路由时需要
	AuthService service.AuthService // 可为空,仅在启用 /admin/auth 时需要
}

// Register 挂载所有路由与全局中间件
func Register(r *gin.Engine, deps Deps) {
	// 全局中间件(顺序:trace → recover → cors → accesslog)
	r.Use(
		middleware.TraceID(),
		middleware.Recover(),
		middleware.CORS(),
		middleware.AccessLog(),
	)

	health := public.NewHealthHandler(deps.ServiceName, deps.Version)

	api := r.Group("/api/v1")
	{
		// ---- 公开接口 ----
		pub := api.Group("/public")
		{
			pub.GET("/health", health.Check)
		}

		// ---- 管理接口 ----
		adm := api.Group("/admin")
		{
			if deps.AuthService != nil {
				authHandler := admin.NewAuthHandler(deps.AuthService)
				adm.POST("/auth/login", authHandler.Login)
			}
			if deps.Signer != nil {
				protected := adm.Group("")
				protected.Use(middleware.JWTAuth(deps.Signer))
				protected.GET("/ping", admin.Ping)
			}
		}
	}
}
