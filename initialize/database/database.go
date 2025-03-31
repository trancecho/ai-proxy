package database

import (
	"fmt"
	"github.com/trancecho/ai-proxy/config"
	"github.com/trancecho/ai-proxy/po"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// db 是全局变量，用来存储数据库连接实例
var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	// 加载配置文件
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	// 从配置中获取数据库相关信息
	dbConfig := cfg.Database

	// 构造连接字符串（DSN）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DatabaseName,
	)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表结构
	err = db.AutoMigrate(&po.RequestLog{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 返回数据库实例
	return db
}

//// GetDB 获取数据库实例
//func GetDB() *gorm.DB {
//	if db == nil {
//		log.Fatal("Database is not initialized. Please call InitDB() first.")
//	}
//	return db
//}
