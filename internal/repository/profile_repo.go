package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yinyin/myblog/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProfileRepo 个人资料仓储(单条记录,id=1)
type ProfileRepo interface {
	Get(ctx context.Context) (*model.Profile, error)
	Upsert(ctx context.Context, p *model.Profile) error
}

type profileRepo struct{ db *gorm.DB }

// NewProfileRepo 构造
func NewProfileRepo(db *gorm.DB) ProfileRepo { return &profileRepo{db: db} }

func (r *profileRepo) Get(ctx context.Context) (*model.Profile, error) {
	var p model.Profile
	if err := r.db.WithContext(ctx).First(&p, 1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

// Upsert 插入或更新唯一的 profile 记录(id=1)。
// 使用 ON DUPLICATE KEY UPDATE,明确更新字段,避免 Save 带上
// 零值的 created_at 在严格模式下报错。
func (r *profileRepo) Upsert(ctx context.Context, p *model.Profile) error {
	p.ID = 1
	now := time.Now()
	p.UpdatedAt = now
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name", "bio", "avatar", "email",
				"github", "twitter", "linkedin", "website",
				"updated_at",
			}),
		}).
		Create(p).Error
	if err != nil {
		return fmt.Errorf("upsert profile: %w", err)
	}
	return nil
}
