package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	ProjectName       string `gorm:"uniqueIndex"`
	VirtualenvName    string
	VirtualenvPath    string
	VirtualenvVersion string
}

func (Project) TableName() string {
	return "project"
}
