package model

// 用户角色
const (
	RoleAdmin  = "admin"
	RoleAuthor = "author"
	RoleReader = "reader"
)

// 用户状态
const (
	UserStatusActive   int8 = 1
	UserStatusDisabled int8 = 0
)

// User 用户表
type User struct {
	BaseModel
	Username     string `gorm:"type:varchar(32);uniqueIndex;not null"  json:"username"`
	Email        string `gorm:"type:varchar(128);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(128);not null"             json:"-"`
	Nickname     string `gorm:"type:varchar(64);not null;default:''"   json:"nickname"`
	Avatar       string `gorm:"type:varchar(255);not null;default:''"  json:"avatar"`
	Bio          string `gorm:"type:varchar(255);not null;default:''"  json:"bio"`
	Role         string `gorm:"type:varchar(16);not null;default:'author'" json:"role"`
	Status       int8   `gorm:"type:tinyint;not null;default:1"        json:"status"`
}

// TableName GORM 表名
func (User) TableName() string { return "users" }

// IsActive 是否启用
func (u *User) IsActive() bool { return u.Status == UserStatusActive }

// IsAdmin 管理员
func (u *User) IsAdmin() bool { return u.Role == RoleAdmin }
