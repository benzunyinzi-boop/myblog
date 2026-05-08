package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/yinyin/myblog/internal/model"
	"gorm.io/gorm"
)

// ProfileRepo 个人资料仓储
type ProfileRepo interface {
	// Get 获取唯一的 profile 记录(不存在返回 ErrNotFound)
	Get(ctx context.Context) (*model.Profile, error)
	// Upsert 更新或插入(profiles 表只有一条记录,id=1)
	Upsert(ctx context.Context, p *model.Profile) error
}

type profileRepo struct{ db *gorm.DB }

// NewProfileRepo 构造
func NewProfileRepo(db *gorm.DB) ProfileRepo { return &profileRepo{db: db} }

func (r *profileRepo) Get(ctx context.Context) (*model.Profile, error) {
	var p model.Profile
	// profiles 表只有一条记录,id=1
	if err := r.db.WithContext(ctx).First(&p, 1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

func (r *profileRepo) Upsert(ctx context.Context, p *model.Profile) error {
	// 强制 id=1,确保只有一条记录
	p.ID = 1
	// Save 会根据主键判断是 INSERT 还是 UPDATE
	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return fmt.Errorf("upsert profile: %w", err)
	}
	return nil
}
