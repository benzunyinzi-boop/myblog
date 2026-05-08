package dto

// ProfileResp 个人资料响应
type ProfileResp struct {
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	GitHub   string `json:"github"`
	Twitter  string `json:"twitter"`
	LinkedIn string `json:"linkedin"`
	Website  string `json:"website"`
}

// ProfileUpdateReq 更新个人资料
type ProfileUpdateReq struct {
	Name     string `json:"name"     binding:"required,max=100"`
	Bio      string `json:"bio"      binding:"required"`
	Avatar   string `json:"avatar"   binding:"max=255"`
	Email    string `json:"email"    binding:"required,email,max=100"`
	GitHub   string `json:"github"   binding:"max=255"`
	Twitter  string `json:"twitter"  binding:"max=255"`
	LinkedIn string `json:"linkedin" binding:"max=255"`
	Website  string `json:"website"  binding:"max=255"`
}
