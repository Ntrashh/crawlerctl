package models

type Program struct {
	Id           int
	Name         string
	ProjectID    uint    // 外键，关联 Project 表的 ID
	Project      Project // 使用 GORM 自动加载关联的 Project
	StartCommand string
}
