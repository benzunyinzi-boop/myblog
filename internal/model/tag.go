package model

// Tag 文章标签
type Tag struct {
	BaseModel
	Name string `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`
	Slug string `gorm:"type:varchar(64);uniqueIndex;not null" json:"slug"`
}

// TableName GORM 表名
func (Tag) TableName() string { return "tags" }

// ArticleTag 文章-标签关联
type ArticleTag struct {
	ArticleID uint64 `gorm:"primaryKey;column:article_id" json:"article_id"`
	TagID     uint64 `gorm:"primaryKey;column:tag_id;index:idx_tag_id" json:"tag_id"`
}

// TableName GORM 表名
func (ArticleTag) TableName() string { return "article_tags" }
