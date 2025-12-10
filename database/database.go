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

// ---------- 获取环境变量，如果不存在则使用默认值 ----------
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

// ---------- 加载数据库配置 ----------
func loadConfig() DBConfig {
	return DBConfig{
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		Name:     getEnv("DB_NAME", "bgs"),
	}
}

// ---------- 创建数据库（如不存在） ----------
func createDatabase(cfg DBConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ 无法连接 MySQL: %v", err)
	}
	defer sqlDB.Close()

	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(5)

	_, err = sqlDB.Exec(
		fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;",
			cfg.Name,
		),
	)
	if err != nil {
		log.Fatalf("❌ 创建数据库失败: %v", err)
	}

	log.Printf("✅ 数据库检查完成：%s 已存在或成功创建\n", cfg.Name)
}

// ---------- 初始化数据库 ----------
func InitDB() {
	cfg := loadConfig()

	// 1. 创建数据库（如果不存在）
	createDatabase(cfg)

	// 2. 使用 GORM 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ GORM 无法连接数据库: %v", err)
	}

	log.Println("✅ GORM 已成功连接数据库")

	// 3. 自动迁移模型（包含所有 models 下的主要表以及多对多 关联表）
	if err := DB.AutoMigrate(
		// &models.Card{},
		&models.Heroes{},
		// &models.Classes{},
		&models.Keywords{},
		// &models.Minions{},
		// &models.Anomalies{},
		&models.Quests{},
		// &models.Rewards{},
		// &models.Spells{},
		&models.Trinkets{},
	); err != nil {
		log.Fatalf("❌ 自动迁移模型表失败: %v", err)
	}

	log.Println("✅ 所有模型表已创建或更新完成")
}
