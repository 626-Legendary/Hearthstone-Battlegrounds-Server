package database

import (
	"bgs-server/models"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

// ---------- 环境变量辅助 ----------
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func loadConfig() DBConfig {
	return DBConfig{
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		Name:     getEnv("DB_NAME", "bgs"),
	}
}

// ---------- 创建数据库（如果不存在） ----------
func createDatabase(cfg DBConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("MySQL 连接失败: %v", err)
	}
	defer sqlDB.Close()

	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(5)

	_, err = sqlDB.Exec(
		"CREATE DATABASE IF NOT EXISTS `" + cfg.Name + "` " +
			"DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;",
	)
	if err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	log.Println("✅ 数据库存在或已创建:", cfg.Name)
}

// ---------- 初始化数据库 ----------
func InitDB() {
	cfg := loadConfig()

	// 1. 确保数据库存在
	createDatabase(cfg)

	// 2. GORM 连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("GORM 连接数据库失败: %v", err)
	}

	log.Println("✅ GORM 连接成功")

	// 3. 自动迁移：只创建 Hero 表
	if err := DB.AutoMigrate(&models.Hero{}); err != nil {
		log.Fatalf("迁移 Hero 表失败: %v", err)
	}

	log.Println("✅ Hero 表结构已自动迁移完成")
}
