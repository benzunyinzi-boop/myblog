package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/yinyin/myblog/internal/dto"
	"github.com/yinyin/myblog/internal/model"
	"github.com/yinyin/myblog/internal/pkg/dberr"
	"github.com/yinyin/myblog/internal/repository"
)

// 分类服务错误(handler 转 errcode)
var (
	ErrCategoryNotFound  = errors.New("service: category not found")
	ErrCategoryDuplicate = errors.New("service: category name or slug duplicated")
)

// CategoryService 分类服务
type CategoryService interface {
	Create(ctx context.Context, req dto.CategoryCreateReq) (*dto.CategoryResp, error)
	Update(ctx context.Context, id uint64, req dto.CategoryUpdateReq) (*dto.CategoryResp, error)
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*dto.CategoryResp, error)
	List(ctx context.Context) ([]*dto.CategoryResp, error)
}

type categoryService struct {
	repo repository.CategoryRepo
}

// NewCategoryService 构造
func NewCategoryService(repo repository.CategoryRepo) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, req dto.CategoryCreateReq) (*dto.CategoryResp, error) {
	c := &model.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		if dberr.IsDuplicateKey(err) {
			return nil, ErrCategoryDuplicate
		}
		return nil, fmt.Errorf("category create: %w", err)
	}
	return toCategoryResp(c), nil
}

func (s *categoryService) Update(ctx context.Context, id uint64, req dto.CategoryUpdateReq) (*dto.CategoryResp, error) {
	c := &model.Category{Name: req.Name, Slug: req.Slug, Description: req.Description, SortOrder: req.SortOrder}
	c.ID = id
	if err := s.repo.Update(ctx, c); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCategoryNotFound
		}
		if dberr.IsDuplicateKey(err) {
			return nil, ErrCategoryDuplicate
		}
		return nil, fmt.Errorf("category update: %w", err)
	}
	// 返回最新数据
	return s.Get(ctx, id)
}

func (s *categoryService) Delete(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrCategoryNotFound
		}
		return fmt.Errorf("category delete: %w", err)
	}
	return nil
}

func (s *categoryService) Get(ctx context.Context, id uint64) (*dto.CategoryResp, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("category get: %w", err)
	}
	return toCategoryResp(c), nil
}

func (s *categoryService) List(ctx context.Context) ([]*dto.CategoryResp, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("category list: %w", err)
	}
	out := make([]*dto.CategoryResp, 0, len(list))
	for _, c := range list {
		out = append(out, toCategoryResp(c))
	}
	return out, nil
}

func toCategoryResp(c *model.Category) *dto.CategoryResp {
	return &dto.CategoryResp{
		ID: c.ID, Name: c.Name, Slug: c.Slug,
		Description: c.Description, SortOrder: c.SortOrder,
	}
}
