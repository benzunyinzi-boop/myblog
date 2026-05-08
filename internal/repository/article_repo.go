package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yinyin/myblog/internal/model"
	"gorm.io/gorm"
)

// ArticleListFilter 文章列表过滤条件
type ArticleListFilter struct {
	Status     string // "" / draft / published
	CategoryID uint64 // 0 表示不过滤
	TagID      uint64 // 0 表示不过滤
	Keyword    string // 在 title/summary 模糊匹配
	AuthorID   uint64 // 0 表示不过滤(admin 场景看自己文章时可用)

	// 分页
	Page     int // 1-based
	PageSize int // 1-100
	// 排序
	OrderBy string // "published_at DESC" / "id DESC",为空走默认
}

// ArticleRepo 文章仓储
type ArticleRepo interface {
	// 事务包装,用于"文章+article_tags"联写
	Tx(ctx context.Context, fn func(txRepo ArticleRepo) error) error

	Create(ctx context.Context, a *model.Article) error
	Update(ctx context.Context, a *model.Article, fields map[string]any) error
	Delete(ctx context.Context, id uint64) error

	GetByID(ctx context.Context, id uint64) (*model.Article, error)
	GetBySlug(ctx context.Context, slug string) (*model.Article, error)
	List(ctx context.Context, f ArticleListFilter) (items []*model.Article, total int64, err error)

	// 浏览计数,走原子 +1
	IncrViewCount(ctx context.Context, id uint64) error

	// ---- article_tags 关联 ----
	ReplaceTags(ctx context.Context, articleID uint64, tagIDs []uint64) error
	GetTagIDs(ctx context.Context, articleID uint64) ([]uint64, error)
}

type articleRepo struct{ db *gorm.DB }

// NewArticleRepo 构造
func NewArticleRepo(db *gorm.DB) ArticleRepo { return &articleRepo{db: db} }

func (r *articleRepo) Tx(ctx context.Context, fn func(ArticleRepo) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&articleRepo{db: tx})
	})
}

func (r *articleRepo) Create(ctx context.Context, a *model.Article) error {
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		return fmt.Errorf("create article: %w", err)
	}
	return nil
}

// Update 按白名单字段更新,避免覆盖 created_at / view_count 等。
// 传入的 fields 为 column→value,所有 key 必须是合法字段名(service 层控制白名单)。
func (r *articleRepo) Update(ctx context.Context, a *model.Article, fields map[string]any) error {
	if a.ID == 0 {
		return fmt.Errorf("update article: id required")
	}
	if len(fields) == 0 {
		return nil
	}
	res := r.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", a.ID).Updates(fields)
	if res.Error != nil {
		return fmt.Errorf("update article: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *articleRepo) Delete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&model.Article{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete article: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *articleRepo) GetByID(ctx context.Context, id uint64) (*model.Article, error) {
	var a model.Article
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get article by id: %w", err)
	}
	return &a, nil
}

func (r *articleRepo) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	var a model.Article
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get article by slug: %w", err)
	}
	return &a, nil
}

func (r *articleRepo) List(ctx context.Context, f ArticleListFilter) ([]*model.Article, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Article{})

	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.CategoryID > 0 {
		q = q.Where("category_id = ?", f.CategoryID)
	}
	if f.AuthorID > 0 {
		q = q.Where("author_id = ?", f.AuthorID)
	}
	if kw := strings.TrimSpace(f.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("title LIKE ? OR summary LIKE ?", like, like)
	}
	if f.TagID > 0 {
		q = q.Where("id IN (SELECT article_id FROM article_tags WHERE tag_id = ?)", f.TagID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count articles: %w", err)
	}
	if total == 0 {
		return nil, 0, nil
	}

	order := f.OrderBy
	if order == "" {
		// 已发布按 published_at 倒序;其他按 id 倒序
		if f.Status == model.ArticleStatusPublished {
			order = "published_at DESC, id DESC"
		} else {
			order = "id DESC"
		}
	}

	page, pageSize := normalizePaging(f.Page, f.PageSize)
	var items []*model.Article
	err := q.Order(order).Limit(pageSize).Offset((page - 1) * pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list articles: %w", err)
	}
	return items, total, nil
}

func (r *articleRepo) IncrViewCount(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	if res.Error != nil {
		return fmt.Errorf("incr view_count: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *articleRepo) ReplaceTags(ctx context.Context, articleID uint64, tagIDs []uint64) error {
	// 先删后插,简单粗暴,交给事务包一层即可
	if err := r.db.WithContext(ctx).
		Where("article_id = ?", articleID).
		Delete(&model.ArticleTag{}).Error; err != nil {
		return fmt.Errorf("clear article_tags: %w", err)
	}
	if len(tagIDs) == 0 {
		return nil
	}
	rows := make([]model.ArticleTag, 0, len(tagIDs))
	seen := make(map[uint64]struct{}, len(tagIDs))
	for _, tid := range tagIDs {
		if tid == 0 {
			continue
		}
		if _, dup := seen[tid]; dup {
			continue
		}
		seen[tid] = struct{}{}
		rows = append(rows, model.ArticleTag{ArticleID: articleID, TagID: tid})
	}
	if len(rows) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(&rows).Error; err != nil {
		return fmt.Errorf("insert article_tags: %w", err)
	}
	return nil
}

func (r *articleRepo) GetTagIDs(ctx context.Context, articleID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).Model(&model.ArticleTag{}).
		Where("article_id = ?", articleID).
		Pluck("tag_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("get tag ids: %w", err)
	}
	return ids, nil
}

func normalizePaging(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return page, size
}
