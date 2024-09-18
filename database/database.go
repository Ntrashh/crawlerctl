package database

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/config"
	"github.com/Ntrashh/crawlerctl/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	DB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/.crawlerctl/crawlerptl.db", config.AppConfig.Path)), &gorm.Config{})
	if err != nil {
		fmt.Printf("无法连接到数据库: %v \n", err)
	}

	// 自动迁移（AutoMigrate）您的模型
	err = DB.AutoMigrate(&models.Project{})
	if err != nil {
		fmt.Printf("数据库迁移失败: %v \n", err)

	}
}
