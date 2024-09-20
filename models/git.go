package models

import "gorm.io/gorm"

type Git struct {
	gorm.Model
	ProjectID int `gorm:"index"`
	GitType   string
	GitPath   string
	UserName  string
	Password  string
	GitBranch string
}
