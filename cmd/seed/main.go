// Package main 简单的种子工具:创建/重置管理员账号。
// Usage:
//
//	go run ./cmd/seed -username admin -password 'Admin@123' [-email ...] [-nickname ...]
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/yinyin/myblog/internal/bootstrap"
	"github.com/yinyin/myblog/internal/config"
	"github.com/yinyin/myblog/internal/model"
	"github.com/yinyin/myblog/internal/pkg/logger"
	"github.com/yinyin/myblog/internal/pkg/password"
	"github.com/yinyin/myblog/internal/repository"
	"go.uber.org/zap"
)

func main() {
	var (
		cfgPath  string
		username string
		pw       string
		email    string
		nickname string
	)
	flag.StringVar(&cfgPath, "config", "", "path to config file")
	flag.StringVar(&username, "username", "admin", "admin username")
	flag.StringVar(&pw, "password", "", "admin password (required)")
	flag.StringVar(&email, "email", "admin@yinyin.dev", "admin email")
	flag.StringVar(&nickname, "nickname", "Admin", "admin nickname")
	flag.Parse()

	if pw == "" {
		fmt.Fprintln(os.Stderr, "missing -password")
		os.Exit(2)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		die("load config", err)
	}
	if err := logger.Init(cfg.Log); err != nil {
		die("init logger", err)
	}
	defer logger.Sync()

	db, err := bootstrap.NewMySQL(cfg.MySQL)
	if err != nil {
		die("connect mysql", err)
	}

	users := repository.NewUserRepo(db)
	ctx := context.Background()

	hash, err := password.Hash(pw)
	if err != nil {
		die("hash password", err)
	}

	existing, err := users.GetByUsername(ctx, username)
	switch {
	case err == nil:
		// 已存在 → 原地更新密码与角色(保留其他字段)
		existing.PasswordHash = hash
		existing.Role = model.RoleAdmin
		existing.Status = model.UserStatusActive
		if err := db.WithContext(ctx).Save(existing).Error; err != nil {
			die("update admin", err)
		}
		logger.L().Info("admin password updated", zap.String("username", username))
	case errors.Is(err, repository.ErrNotFound):
		u := &model.User{
			Username:     username,
			Email:        email,
			PasswordHash: hash,
			Nickname:     nickname,
			Role:         model.RoleAdmin,
			Status:       model.UserStatusActive,
		}
		if err := users.Create(ctx, u); err != nil {
			die("create admin", err)
		}
		logger.L().Info("admin created",
			zap.String("username", username),
			zap.Uint64("id", u.ID),
		)
	default:
		die("lookup admin", err)
	}
}

func die(stage string, err error) {
	fmt.Fprintf(os.Stderr, "seed: %s: %v\n", stage, err)
	os.Exit(1)
}
