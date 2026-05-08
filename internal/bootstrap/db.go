package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/yinyin/myblog/internal/config"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// zapWriter 把 GORM 日志桥接到 zap(Info 级)
type zapWriter struct{}

func (zapWriter) Printf(format string, args ...any) {
	logger.L().Info(fmt.Sprintf(format, args...))
}

// NewMySQL 建立 GORM 连接并做健康检查。
func NewMySQL(cfg config.MySQLConfig) (*gorm.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("mysql.dsn is empty")
	}

	gormLog := gormlogger.New(
		zapWriter{},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger:                                   gormLog,
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc:                                  time.Now,
	})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get *sql.DB: %w", err)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	logger.L().Info("mysql connected",
		zap.Int("max_open_conns", cfg.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
	)
	return db, nil
}
