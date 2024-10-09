package database

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/config"
	"github.com/Ntrashh/crawlerctl/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDatabase() *gorm.DB {
	var err error

	DB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/.crawlerctl/crawlerptl.db", config.AppConfig.Path)), &gorm.Config{})
	if err != nil {

		fmt.Printf("无法连接到数据库: %v \n", err)
	}
	// 执行迁移
	err = DB.AutoMigrate(&models.Project{}, &models.Git{}, &models.Program{})
	if err != nil {
		log.Fatal("迁移失败:", err)
	}
	return DB
}
