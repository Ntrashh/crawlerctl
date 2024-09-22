package models

import (
	"gorm.io/gorm"
)

type Project struct {
	ProjectName       string `gorm:"uniqueIndex"`
	VirtualenvName    string
	VirtualenvPath    string
	VirtualenvVersion string
	SavePath          string
	gorm.Model
}

func (Project) TableName() string {
	return "project"
}
