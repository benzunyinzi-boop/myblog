package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yinyin/myblog/internal/config"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"go.uber.org/zap"
)

// NewRedis 建立 Redis 客户端并健康检查。
func NewRedis(cfg config.RedisConfig) (*redis.Client, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("redis.addr is empty")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	logger.L().Info("redis connected", zap.String("addr", cfg.Addr), zap.Int("db", cfg.DB))
	return rdb, nil
}
