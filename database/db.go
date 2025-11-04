package database

import (
	"log"

	"github.com/Arterning/go-rag/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("documents.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(&models.Document{}, &models.DocumentChunk{})
	if err != nil {
		return err
	}

	log.Println("数据库初始化成功")
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}
