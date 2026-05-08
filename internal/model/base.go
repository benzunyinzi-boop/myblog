// Package model 数据库模型(GORM)。
package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 所有业务表统一嵌入(CLAUDE.md 规范)
type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"        json:"id"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime:milli"   json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime:milli"   json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                           json:"-"`
}
