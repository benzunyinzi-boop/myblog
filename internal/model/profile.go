package model

// Profile 个人资料(关于我)
type Profile struct {
	BaseModel
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Bio         string `gorm:"type:text;not null"         json:"bio"`
	Avatar      string `gorm:"type:varchar(255);not null" json:"avatar"`
	Email       string `gorm:"type:varchar(100);not null" json:"email"`
	GitHub      string `gorm:"type:varchar(255);not null" json:"github"`
	Twitter     string `gorm:"type:varchar(255);not null" json:"twitter"`
	LinkedIn    string `gorm:"type:varchar(255);not null" json:"linkedin"`
	Website     string `gorm:"type:varchar(255);not null" json:"website"`
}

// TableName GORM 表名
func (Profile) TableName() string { return "profiles" }
