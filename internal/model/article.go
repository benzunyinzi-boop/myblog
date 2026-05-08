package model

// 文章状态
const (
	ArticleStatusDraft     = "draft"     // 草稿
	ArticleStatusPublished = "published" // 已发布
)

// Article 文章
type Article struct {
	BaseModel
	Title       string  `gorm:"type:varchar(200);not null"                json:"title"`
	Slug        string  `gorm:"type:varchar(200);uniqueIndex;not null"    json:"slug"`
	Summary     string  `gorm:"type:varchar(500);not null;default:''"     json:"summary"`
	Content     string  `gorm:"type:mediumtext;not null"                  json:"content"`
	CoverImage  string  `gorm:"type:varchar(255);not null;default:''"     json:"cover_image"`
	CategoryID  uint64  `gorm:"type:bigint unsigned;not null;default:0;index:idx_category_id" json:"category_id"`
	AuthorID    uint64  `gorm:"type:bigint unsigned;not null;index:idx_author_id" json:"author_id"`
	Status      string  `gorm:"type:varchar(16);not null;default:'draft';index:idx_status_published_at" json:"status"`
	ViewCount   int     `gorm:"type:int;not null;default:0"               json:"view_count"`
	PublishedAt *int64  `gorm:"type:bigint;index:idx_status_published_at" json:"published_at,omitempty"` // unix 秒,允许 null
}

// TableName GORM 表名
func (Article) TableName() string { return "articles" }

// IsPublished 是否已发布
func (a *Article) IsPublished() bool { return a.Status == ArticleStatusPublished }

// IsDraft 是否草稿
func (a *Article) IsDraft() bool { return a.Status == ArticleStatusDraft }
