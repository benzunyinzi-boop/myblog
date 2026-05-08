package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yinyin/myblog/internal/dto"
	"github.com/yinyin/myblog/internal/model"
	"github.com/yinyin/myblog/internal/pkg/dberr"
	"github.com/yinyin/myblog/internal/repository"
)

// ArticleService 文章服务
type ArticleService interface {
	Create(ctx context.Context, authorID uint64, req dto.ArticleCreateReq) (*dto.ArticleDetailResp, error)
	Update(ctx context.Context, id uint64, req dto.ArticleUpdateReq) (*dto.ArticleDetailResp, error)
	Delete(ctx context.Context, id uint64) error

	GetByID(ctx context.Context, id uint64) (*dto.ArticleDetailResp, error)
	GetBySlug(ctx context.Context, slug string) (*dto.ArticleDetailResp, error)
	List(ctx context.Context, req dto.ArticleListReq) (*dto.ArticleListResp, error)

	// Publish 将草稿转为已发布(已发布时返回 ErrArticleConflict)
	Publish(ctx context.Context, id uint64) (*dto.ArticleDetailResp, error)
	// Unpublish 将已发布转为草稿(草稿时返回 ErrArticleConflict)
	Unpublish(ctx context.Context, id uint64) (*dto.ArticleDetailResp, error)

	// IncrView 浏览计数 +1(不存在返回 ErrArticleNotFound)
	IncrView(ctx context.Context, id uint64) error
}

type articleService struct {
	articles   repository.ArticleRepo
	categories repository.CategoryRepo
	tags       repository.TagRepo
}

// NewArticleService 构造
func NewArticleService(a repository.ArticleRepo, c repository.CategoryRepo, t repository.TagRepo) ArticleService {
	return &articleService{articles: a, categories: c, tags: t}
}

func (s *articleService) Create(ctx context.Context, authorID uint64, req dto.ArticleCreateReq) (*dto.ArticleDetailResp, error) {
	if err := s.validateCategory(ctx, req.CategoryID); err != nil {
		return nil, err
	}
	if err := s.validateTagIDs(ctx, req.TagIDs); err != nil {
		return nil, err
	}

	status := req.Status
	if status == "" {
		status = model.ArticleStatusDraft
	}

	a := &model.Article{
		Title:      req.Title,
		Slug:       req.Slug,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		CategoryID: req.CategoryID,
		AuthorID:   authorID,
		Status:     status,
	}
	if status == model.ArticleStatusPublished {
		now := time.Now().Unix()
		a.PublishedAt = &now
	}

	err := s.articles.Tx(ctx, func(tx repository.ArticleRepo) error {
		if err := tx.Create(ctx, a); err != nil {
			if dberr.IsDuplicateKey(err) {
				return ErrArticleDuplicate
			}
			return fmt.Errorf("article create: %w", err)
		}
		if err := tx.ReplaceTags(ctx, a.ID, req.TagIDs); err != nil {
			return fmt.Errorf("article tags: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.buildDetail(ctx, a)
}

func (s *articleService) Update(ctx context.Context, id uint64, req dto.ArticleUpdateReq) (*dto.ArticleDetailResp, error) {
	if err := s.validateCategory(ctx, req.CategoryID); err != nil {
		return nil, err
	}
	if err := s.validateTagIDs(ctx, req.TagIDs); err != nil {
		return nil, err
	}

	existed, err := s.articles.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, fmt.Errorf("article update get: %w", err)
	}

	fields := map[string]any{
		"title":       req.Title,
		"slug":        req.Slug,
		"summary":     req.Summary,
		"content":     req.Content,
		"cover_image": req.CoverImage,
		"category_id": req.CategoryID,
	}

	err = s.articles.Tx(ctx, func(tx repository.ArticleRepo) error {
		if err := tx.Update(ctx, existed, fields); err != nil {
			if dberr.IsDuplicateKey(err) {
				return ErrArticleDuplicate
			}
			if errors.Is(err, repository.ErrNotFound) {
				return ErrArticleNotFound
			}
			return fmt.Errorf("article update: %w", err)
		}
		if err := tx.ReplaceTags(ctx, id, req.TagIDs); err != nil {
			return fmt.Errorf("article tags: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *articleService) Delete(ctx context.Context, id uint64) error {
	err := s.articles.Tx(ctx, func(tx repository.ArticleRepo) error {
		// 先清关联,再软删文章(article_tags 没有 soft delete)
		if err := tx.ReplaceTags(ctx, id, nil); err != nil {
			return fmt.Errorf("article delete tags: %w", err)
		}
		if err := tx.Delete(ctx, id); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrArticleNotFound
			}
			return fmt.Errorf("article delete: %w", err)
		}
		return nil
	})
	return err
}

func (s *articleService) GetByID(ctx context.Context, id uint64) (*dto.ArticleDetailResp, error) {
	a, err := s.articles.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, fmt.Errorf("article get: %w", err)
	}
	return s.buildDetail(ctx, a)
}

func (s *articleService) GetBySlug(ctx context.Context, slug string) (*dto.ArticleDetailResp, error) {
	a, err := s.articles.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, fmt.Errorf("article get by slug: %w", err)
	}
	return s.buildDetail(ctx, a)
}

func (s *articleService) List(ctx context.Context, req dto.ArticleListReq) (*dto.ArticleListResp, error) {
	items, total, err := s.articles.List(ctx, repository.ArticleListFilter{
		Status:     req.Status,
		CategoryID: req.CategoryID,
		TagID:      req.TagID,
		Keyword:    req.Keyword,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("article list: %w", err)
	}

	summaries := make([]*dto.ArticleSummaryResp, 0, len(items))
	for _, a := range items {
		sm, err := s.buildSummary(ctx, a)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, sm)
	}

	page, size := req.Page, req.PageSize
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	return &dto.ArticleListResp{Items: summaries, Total: total, Page: page, PageSize: size}, nil
}

func (s *articleService) Publish(ctx context.Context, id uint64) (*dto.ArticleDetailResp, error) {
	a, err := s.articles.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, fmt.Errorf("article publish get: %w", err)
	}
	if a.IsPublished() {
		return nil, ErrArticleConflict
	}
	now := time.Now().Unix()
	if err := s.articles.Update(ctx, a, map[string]any{
		"status":       model.ArticleStatusPublished,
		"published_at": now,
	}); err != nil {
		return nil, fmt.Errorf("article publish update: %w", err)
	}
	return s.GetByID(ctx, id)
}

func (s *articleService) Unpublish(ctx context.Context, id uint64) (*dto.ArticleDetailResp, error) {
	a, err := s.articles.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, fmt.Errorf("article unpublish get: %w", err)
	}
	if a.IsDraft() {
		return nil, ErrArticleConflict
	}
	if err := s.articles.Update(ctx, a, map[string]any{
		"status": model.ArticleStatusDraft,
	}); err != nil {
		return nil, fmt.Errorf("article unpublish update: %w", err)
	}
	return s.GetByID(ctx, id)
}

func (s *articleService) IncrView(ctx context.Context, id uint64) error {
	if err := s.articles.IncrViewCount(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrArticleNotFound
		}
		return fmt.Errorf("article incr view: %w", err)
	}
	return nil
}

// ---- helpers ----

func (s *articleService) validateCategory(ctx context.Context, cid uint64) error {
	if cid == 0 {
		return nil // 允许不选分类
	}
	if _, err := s.categories.GetByID(ctx, cid); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrCategoryNotFound
		}
		return fmt.Errorf("validate category: %w", err)
	}
	return nil
}

func (s *articleService) validateTagIDs(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	existing, err := s.tags.ListByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("validate tags: %w", err)
	}
	// 去重比较
	set := make(map[uint64]struct{}, len(existing))
	for _, t := range existing {
		set[t.ID] = struct{}{}
	}
	for _, id := range ids {
		if _, ok := set[id]; !ok {
			return ErrTagNotFound
		}
	}
	return nil
}

func (s *articleService) buildSummary(ctx context.Context, a *model.Article) (*dto.ArticleSummaryResp, error) {
	tagIDs, err := s.articles.GetTagIDs(ctx, a.ID)
	if err != nil {
		return nil, err
	}
	tags, err := s.tags.ListByIDs(ctx, tagIDs)
	if err != nil {
		return nil, err
	}
	tagResps := make([]dto.TagResp, 0, len(tags))
	for _, t := range tags {
		tagResps = append(tagResps, dto.TagResp{ID: t.ID, Name: t.Name, Slug: t.Slug})
	}
	return &dto.ArticleSummaryResp{
		ID:          a.ID,
		Title:       a.Title,
		Slug:        a.Slug,
		Summary:     a.Summary,
		CoverImage:  a.CoverImage,
		CategoryID:  a.CategoryID,
		AuthorID:    a.AuthorID,
		Status:      a.Status,
		ViewCount:   a.ViewCount,
		PublishedAt: a.PublishedAt,
		CreatedAt:   a.CreatedAt.Unix(),
		Tags:        tagResps,
	}, nil
}

func (s *articleService) buildDetail(ctx context.Context, a *model.Article) (*dto.ArticleDetailResp, error) {
	sm, err := s.buildSummary(ctx, a)
	if err != nil {
		return nil, err
	}
	return &dto.ArticleDetailResp{ArticleSummaryResp: *sm, Content: a.Content}, nil
}
