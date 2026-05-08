package dto

// CategoryCreateReq 创建分类
type CategoryCreateReq struct {
	Name        string `json:"name"        binding:"required,max=64"`
	Slug        string `json:"slug"        binding:"required,max=64,slug"`
	Description string `json:"description" binding:"max=255"`
	SortOrder   int    `json:"sort_order"  binding:"gte=-1000,lte=1000"`
}

// CategoryUpdateReq 更新分类
type CategoryUpdateReq struct {
	Name        string `json:"name"        binding:"required,max=64"`
	Slug        string `json:"slug"        binding:"required,max=64,slug"`
	Description string `json:"description" binding:"max=255"`
	SortOrder   int    `json:"sort_order"  binding:"gte=-1000,lte=1000"`
}

// CategoryResp 分类响应
type CategoryResp struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}
