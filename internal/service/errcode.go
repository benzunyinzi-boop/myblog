// Package service 业务逻辑层。
// errcode.go 统一定义 service 层的业务错误,handler 层按类型映射到 pkg/errcode。
package service

import "errors"

// 通用错误
var (
	ErrNotFound  = errors.New("service: resource not found")
	ErrDuplicate = errors.New("service: resource duplicated")
	ErrConflict  = errors.New("service: operation conflict")
)

// 用户相关
var (
	ErrUserNotFound    = errors.New("service: user not found")
	ErrInvalidPassword = errors.New("service: invalid password")
	ErrUserDisabled    = errors.New("service: user disabled")
)

// 分类相关
var (
	ErrCategoryNotFound  = errors.New("service: category not found")
	ErrCategoryDuplicate = errors.New("service: category name or slug duplicated")
	ErrCategoryInUse     = errors.New("service: category in use by articles")
)

// 标签相关
var (
	ErrTagNotFound  = errors.New("service: tag not found")
	ErrTagDuplicate = errors.New("service: tag duplicated")
)

// 文章相关(预留,M3 part 2 使用)
var (
	ErrArticleNotFound  = errors.New("service: article not found")
	ErrArticleDuplicate = errors.New("service: article slug duplicated")
	ErrArticleConflict  = errors.New("service: article status conflict")
)
