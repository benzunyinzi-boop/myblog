package model

// Profile 个人资料(关于我)。profiles 表只保留一条记录(id=1)。
type Profile struct {
	BaseModel
	Name     string `gorm:"column:name;type:varchar(100);not null;default:''"     json:"name"`
	Bio      string `gorm:"column:bio;type:text;not null"                         json:"bio"`
	Avatar   string `gorm:"column:avatar;type:varchar(255);not null;default:''"   json:"avatar"`
	Email    string `gorm:"column:email;type:varchar(100);not null;default:''"    json:"email"`
	GitHub   string `gorm:"column:github;type:varchar(255);not null;default:''"   json:"github"`
	Twitter  string `gorm:"column:twitter;type:varchar(255);not null;default:''"  json:"twitter"`
	LinkedIn string `gorm:"column:linkedin;type:varchar(255);not null;default:''" json:"linkedin"`
	Website  string `gorm:"column:website;type:varchar(255);not null;default:''"  json:"website"`
}

// TableName GORM 表名
func (Profile) TableName() string { return "profiles" }
