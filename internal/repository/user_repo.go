// Package repository 数据访问层。
// 约束:不包含业务判断,只做 CRUD 与查询组合。
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/yinyin/myblog/internal/model"
	"gorm.io/gorm"
)

// ErrNotFound 数据不存在(统一错误,service 层按需翻译成 errcode)
var ErrNotFound = errors.New("repository: record not found")

// UserRepo 用户仓储接口。
type UserRepo interface {
	Create(ctx context.Context, u *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

// userRepo GORM 实现
type userRepo struct {
	db *gorm.DB
}

// NewUserRepo 构造
func NewUserRepo(db *gorm.DB) UserRepo { return &userRepo{db: db} }

func (r *userRepo) Create(ctx context.Context, u *model.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return &u, nil
}
