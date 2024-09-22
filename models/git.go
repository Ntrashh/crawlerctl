package models

import "gorm.io/gorm"

type Git struct {
	gorm.Model
	ProjectID int `gorm:"index"`
	GitPath   string
	UserName  string
	Password  string `gorm:"type:varchar(255);not null" json:"-"`
}

func (Git) TableName() string {
	return "git"
}
