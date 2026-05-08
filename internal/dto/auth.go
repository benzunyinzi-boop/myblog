// Package dto Data Transfer Objects:Req / Resp
package dto

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

// LoginResp 登录响应
type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"` // access token 到期 unix 秒
	User         UserInfoResp `json:"user"`
}

// UserInfoResp 用户基础信息
type UserInfoResp struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}
