package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/yinyin/myblog/internal/model"
	"gorm.io/gorm"
)

// CategoryRepo 分类仓储
type CategoryRepo interface {
	Create(ctx context.Context, c *model.Category) error
	Update(ctx context.Context, c *model.Category) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.Category, error)
	GetBySlug(ctx context.Context, slug string) (*model.Category, error)
	List(ctx context.Context) ([]*model.Category, error)
}

type categoryRepo struct{ db *gorm.DB }

// NewCategoryRepo 构造
func NewCategoryRepo(db *gorm.DB) CategoryRepo { return &categoryRepo{db: db} }

func (r *categoryRepo) Create(ctx context.Context, c *model.Category) error {
	if err := r.db.WithContext(ctx).Create(c).Error; err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (r *categoryRepo) Update(ctx context.Context, c *model.Category) error {
	// 只更新白名单字段,避免覆盖 created_at 等
	fields := map[string]any{
		"name":        c.Name,
		"slug":        c.Slug,
		"description": c.Description,
		"sort_order":  c.SortOrder,
	}
	res := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", c.ID).Updates(fields)
	if res.Error != nil {
		return fmt.Errorf("update category: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *categoryRepo) Delete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&model.Category{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete category: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *categoryRepo) GetByID(ctx context.Context, id uint64) (*model.Category, error) {
	var c model.Category
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get category by id: %w", err)
	}
	return &c, nil
}

func (r *categoryRepo) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	var c model.Category
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get category by slug: %w", err)
	}
	return &c, nil
}

func (r *categoryRepo) List(ctx context.Context) ([]*model.Category, error) {
	var list []*model.Category
	if err := r.db.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return list, nil
}
