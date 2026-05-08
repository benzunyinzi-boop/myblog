package model

// Category 文章分类
type Category struct {
	BaseModel
	Name        string `gorm:"type:varchar(64);uniqueIndex;not null"    json:"name"`
	Slug        string `gorm:"type:varchar(64);uniqueIndex;not null"    json:"slug"`
	Description string `gorm:"type:varchar(255);not null;default:''"    json:"description"`
	SortOrder   int    `gorm:"type:int;not null;default:0;index:idx_sort_order" json:"sort_order"`
}

// TableName GORM 表名
func (Category) TableName() string { return "categories" }
