package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	ProjectName       string `gorm:"uniqueIndex"`
	VirtualEnvName    string
	VirtualEnvPath    string
	VirtualEnvVersion string
}

func (Project) TableName() string {
	return "project"
}
