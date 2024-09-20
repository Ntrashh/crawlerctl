package database

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() *gorm.DB {
	var err error

	DB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/.crawlerctl/crawlerptl.db", config.AppConfig.Path)), &gorm.Config{})
	if err != nil {
		fmt.Printf("无法连接到数据库: %v \n", err)
	}
	return DB

}
