package dto

// ---- 管理端请求 ----

// ArticleCreateReq 创建文章
type ArticleCreateReq struct {
	Title      string   `json:"title"       binding:"required,max=200"`
	Slug       string   `json:"slug"        binding:"required,max=200,slug"`
	Summary    string   `json:"summary"     binding:"max=500"`
	Content    string   `json:"content"     binding:"required"`
	CoverImage string   `json:"cover_image" binding:"max=255"`
	CategoryID uint64   `json:"category_id"`
	TagIDs     []uint64 `json:"tag_ids"     binding:"max=20,dive,gt=0"`
	Status     string   `json:"status"      binding:"omitempty,oneof=draft published"`
}

// ArticleUpdateReq 更新文章(字段与 Create 相同,status 走单独 publish 接口切换)
type ArticleUpdateReq struct {
	Title      string   `json:"title"       binding:"required,max=200"`
	Slug       string   `json:"slug"        binding:"required,max=200,slug"`
	Summary    string   `json:"summary"     binding:"max=500"`
	Content    string   `json:"content"     binding:"required"`
	CoverImage string   `json:"cover_image" binding:"max=255"`
	CategoryID uint64   `json:"category_id"`
	TagIDs     []uint64 `json:"tag_ids"     binding:"max=20,dive,gt=0"`
}

// ArticleListReq 列表查询参数(admin & public 共用,public 强制 status=published)
type ArticleListReq struct {
	Page       int    `form:"page"        binding:"omitempty,gte=1"`
	PageSize   int    `form:"page_size"   binding:"omitempty,gte=1,lte=100"`
	Status     string `form:"status"      binding:"omitempty,oneof=draft published"`
	CategoryID uint64 `form:"category_id"`
	TagID      uint64 `form:"tag_id"`
	Keyword    string `form:"keyword"     binding:"max=100"`
}

// ---- 响应 ----

// ArticleSummaryResp 文章摘要(列表项)
type ArticleSummaryResp struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Summary     string    `json:"summary"`
	CoverImage  string    `json:"cover_image"`
	CategoryID  uint64    `json:"category_id"`
	AuthorID    uint64    `json:"author_id"`
	Status      string    `json:"status"`
	ViewCount   int       `json:"view_count"`
	PublishedAt *int64    `json:"published_at,omitempty"`
	CreatedAt   int64     `json:"created_at"`
	Tags        []TagResp `json:"tags"`
}

// ArticleDetailResp 文章详情(带 content)
type ArticleDetailResp struct {
	ArticleSummaryResp
	Content string `json:"content"`
}

// ArticleListResp 分页列表
type ArticleListResp struct {
	Items    []*ArticleSummaryResp `json:"items"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
}
