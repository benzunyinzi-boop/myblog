package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/yinyin/myblog/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TagRepo 标签仓储
type TagRepo interface {
	Create(ctx context.Context, t *model.Tag) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*model.Tag, error)
	List(ctx context.Context) ([]*model.Tag, error)
	ListByIDs(ctx context.Context, ids []uint64) ([]*model.Tag, error)
	// FirstOrCreateByName 幂等获取/创建(编辑文章时常用)
	FirstOrCreateByName(ctx context.Context, name, slug string) (*model.Tag, error)
}

type tagRepo struct{ db *gorm.DB }

// NewTagRepo 构造
func NewTagRepo(db *gorm.DB) TagRepo { return &tagRepo{db: db} }

func (r *tagRepo) Create(ctx context.Context, t *model.Tag) error {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	return nil
}

func (r *tagRepo) Delete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&model.Tag{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete tag: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *tagRepo) GetByID(ctx context.Context, id uint64) (*model.Tag, error) {
	var t model.Tag
	if err := r.db.WithContext(ctx).First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get tag by id: %w", err)
	}
	return &t, nil
}

func (r *tagRepo) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	var t model.Tag
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get tag by slug: %w", err)
	}
	return &t, nil
}

func (r *tagRepo) List(ctx context.Context) ([]*model.Tag, error) {
	var list []*model.Tag
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	return list, nil
}

func (r *tagRepo) ListByIDs(ctx context.Context, ids []uint64) ([]*model.Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []*model.Tag
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list tags by ids: %w", err)
	}
	return list, nil
}

func (r *tagRepo) FirstOrCreateByName(ctx context.Context, name, slug string) (*model.Tag, error) {
	t := &model.Tag{Name: name, Slug: slug}
	// 唯一索引冲突时不抛错,走 on conflict do nothing,再取一次
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(t).Error
	if err != nil {
		return nil, fmt.Errorf("first or create tag: %w", err)
	}
	if t.ID == 0 {
		// 插入被忽略,说明已存在,按 name 查一次
		return r.getByName(ctx, name)
	}
	return t, nil
}

func (r *tagRepo) getByName(ctx context.Context, name string) (*model.Tag, error) {
	var t model.Tag
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get tag by name: %w", err)
	}
	return &t, nil
}
