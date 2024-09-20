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

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&Project{})
	if err != nil {
		return
	}
}
