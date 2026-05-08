// Package bootstrap 应用启动装配。
package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yinyin/myblog/internal/config"
	myjwt "github.com/yinyin/myblog/internal/pkg/jwt"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/uploader"
	"github.com/yinyin/myblog/internal/pkg/validation"
	"github.com/yinyin/myblog/internal/repository"
	"github.com/yinyin/myblog/internal/router"
	"github.com/yinyin/myblog/internal/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ServiceName 服务名(日志/health 里用)
const ServiceName = "myblog"

// BuildEngine 构建 gin.Engine(供集成测试复用,不起 HTTP 监听)。
// deps 允许为零值,零值时只注册不依赖 DB/Redis/JWT 的路由(如 /public/health)。
func BuildEngine(cfg *config.Config, version string, deps router.Deps) *gin.Engine {
	switch cfg.Server.Mode {
	case gin.ReleaseMode, gin.DebugMode, gin.TestMode:
		gin.SetMode(cfg.Server.Mode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 注册项目级自定义 validator(slug 等)
	if err := validation.Register(); err != nil {
		logger.L().Warn("register validation failed", zap.Error(err))
	}

	deps.ServiceName = ServiceName
	deps.Version = version

	r := gin.New()
	router.Register(r, deps)
	return r
}

// Run 启动服务并阻塞,收到 SIGINT/SIGTERM 优雅关停
func Run(cfgPath, version string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := logger.Init(cfg.Log); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer logger.Sync()

	// ---- 装配依赖 ----
	db, err := NewMySQL(cfg.MySQL)
	if err != nil {
		return fmt.Errorf("init mysql: %w", err)
	}

	rdb, err := NewRedis(cfg.Redis)
	if err != nil {
		return fmt.Errorf("init redis: %w", err)
	}

	signer, err := myjwt.NewSigner(myjwt.Config{
		Secret:     cfg.JWT.Secret,
		AccessTTL:  cfg.JWT.AccessTTL,
		RefreshTTL: cfg.JWT.RefreshTTL,
		Issuer:     cfg.JWT.Issuer,
	})
	if err != nil {
		return fmt.Errorf("init jwt: %w", err)
	}

	userRepo := repository.NewUserRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	tagRepo := repository.NewTagRepo(db)
	articleRepo := repository.NewArticleRepo(db)
	profileRepo := repository.NewProfileRepo(db)
	authSvc := service.NewAuthService(userRepo, signer)
	categorySvc := service.NewCategoryService(categoryRepo)
	tagSvc := service.NewTagService(tagRepo)
	articleSvc := service.NewArticleService(articleRepo, categoryRepo, tagRepo)
	profileSvc := service.NewProfileService(profileRepo)

	// 文件上传器(本地实现)
	fileUploader, err := uploader.NewLocalUploader("./uploads", "/uploads")
	if err != nil {
		logger.L().Fatal("init uploader", zap.Error(err))
	}

	engine := BuildEngine(cfg, version, router.Deps{
		Signer:      signer,
		AuthService: authSvc,
		CategorySvc: categorySvc,
		TagSvc:      tagSvc,
		ArticleSvc:  articleSvc,
		ProfileSvc:  profileSvc,
		Uploader:    fileUploader,
	})

	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动 HTTP
	errCh := make(chan error, 1)
	go func() {
		logger.L().Info("http server starting",
			zap.String("addr", cfg.Server.Addr()),
			zap.String("mode", cfg.Server.Mode),
			zap.String("version", version),
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// 等待信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("http serve: %w", err)
	case sig := <-quit:
		logger.L().Info("shutdown signal", zap.String("signal", sig.String()))
	}

	// 优雅关停
	timeout := cfg.Server.ShutdownTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	// 关闭下游连接
	closeRedis(rdb)
	closeDB(db)

	logger.L().Info("http server stopped gracefully")
	return nil
}

func closeRedis(rdb *redis.Client) {
	if rdb == nil {
		return
	}
	if err := rdb.Close(); err != nil {
		logger.L().Warn("close redis", zap.Error(err))
	}
}

func closeDB(db *gorm.DB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.L().Warn("get *sql.DB for close", zap.Error(err))
		return
	}
	if err := sqlDB.Close(); err != nil {
		logger.L().Warn("close mysql", zap.Error(err))
	}
}
