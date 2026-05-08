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
	Signer      *myjwt.Signer           // 可为空,仅在启用鉴权路由时需要
	AuthService service.AuthService     // 可为空,仅在启用 /admin/auth 时需要
	CategorySvc service.CategoryService // 可为空,仅在启用 category 路由时需要
	TagSvc      service.TagService      // 可为空,仅在启用 tag 路由时需要
	ArticleSvc  service.ArticleService  // 可为空,仅在启用 article 路由时需要
}

// Register 挂载所有路由与全局中间件
func Register(r *gin.Engine, deps Deps) {
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
			if deps.CategorySvc != nil {
				ch := public.NewCategoryHandler(deps.CategorySvc)
				pub.GET("/categories", ch.List)
			}
			if deps.TagSvc != nil {
				th := public.NewTagHandler(deps.TagSvc)
				pub.GET("/tags", th.List)
			}
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

				if deps.CategorySvc != nil {
					ch := admin.NewCategoryHandler(deps.CategorySvc)
					protected.GET("/categories", ch.List)
					protected.POST("/categories", ch.Create)
					protected.PUT("/categories/:id", ch.Update)
					protected.DELETE("/categories/:id", ch.Delete)
				}
				if deps.TagSvc != nil {
					th := admin.NewTagHandler(deps.TagSvc)
					protected.GET("/tags", th.List)
					protected.POST("/tags", th.Create)
					protected.DELETE("/tags/:id", th.Delete)
				}
				if deps.ArticleSvc != nil {
					ah := admin.NewArticleHandler(deps.ArticleSvc)
					protected.GET("/articles", ah.List)
					protected.POST("/articles", ah.Create)
					protected.GET("/articles/:id", ah.Get)
					protected.PUT("/articles/:id", ah.Update)
					protected.DELETE("/articles/:id", ah.Delete)
					protected.POST("/articles/:id/publish", ah.Publish)
					protected.POST("/articles/:id/unpublish", ah.Unpublish)
				}
			}
		}
	}
}
