package dto

// TagCreateReq 创建标签
type TagCreateReq struct {
	Name string `json:"name" binding:"required,max=64"`
	Slug string `json:"slug" binding:"required,max=64,slug"`
}

// TagResp 标签响应
type TagResp struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
