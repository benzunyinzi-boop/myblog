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

// TagService 标签服务
type TagService interface {
	Create(ctx context.Context, req dto.TagCreateReq) (*dto.TagResp, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context) ([]*dto.TagResp, error)
}

type tagService struct{ repo repository.TagRepo }

// NewTagService 构造
func NewTagService(repo repository.TagRepo) TagService { return &tagService{repo: repo} }

func (s *tagService) Create(ctx context.Context, req dto.TagCreateReq) (*dto.TagResp, error) {
	t := &model.Tag{Name: req.Name, Slug: req.Slug}
	if err := s.repo.Create(ctx, t); err != nil {
		if dberr.IsDuplicateKey(err) {
			return nil, ErrTagDuplicate
		}
		return nil, fmt.Errorf("tag create: %w", err)
	}
	return toTagResp(t), nil
}

func (s *tagService) Delete(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTagNotFound
		}
		return fmt.Errorf("tag delete: %w", err)
	}
	return nil
}

func (s *tagService) List(ctx context.Context) ([]*dto.TagResp, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("tag list: %w", err)
	}
	out := make([]*dto.TagResp, 0, len(list))
	for _, t := range list {
		out = append(out, toTagResp(t))
	}
	return out, nil
}

func toTagResp(t *model.Tag) *dto.TagResp {
	return &dto.TagResp{ID: t.ID, Name: t.Name, Slug: t.Slug}
}
