package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接并自动迁移。
//   driver:     sqlite | mysql | postgres
//   dsn:        mysql/postgres 连接串（driver 非 sqlite 时必填）
//   sqlitePath: sqlite 文件路径（仅 sqlite 生效，空则使用内存数据库）
func InitDB(driver, dsn, sqlitePath string) {
	var err error
	DB, err = gorm.Open(openDialector(driver, dsn, sqlitePath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库模型（模型层无数据库专属语法，跨库通用）
	if err = DB.AutoMigrate(&model.User{}, &model.Cluster{}, &model.AuditLog{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 首次启动若无管理员则创建默认账户
	var adminCount int64
	DB.Model(&model.User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}
		if err := DB.Create(&model.User{
			Username: "admin",
			Email:    "admin@example.com",
			Role:     "admin",
			Password: string(hashedPassword),
		}).Error; err != nil {
			log.Fatal("Failed to create admin user:", err)
		}
		log.Println("Created default admin user: admin/admin123")
	}

	log.Printf("Database initialized successfully (driver=%s)", driver)
}

// openDialector 按 driver 选择 GORM Dialector。
func openDialector(driver, dsn, sqlitePath string) gorm.Dialector {
	switch driver {
	case "mysql":
		return mysql.Open(dsn)
	case "postgres", "postgresql":
		return postgres.Open(dsn)
	case "sqlite", "":
		if sqlitePath == "" {
			sqlitePath = "file::memory:?cache=shared"
		} else if dir := filepath.Dir(sqlitePath); dir != "" {
			// 确保数据目录存在
			if err := os.MkdirAll(dir, 0o755); err != nil {
				log.Fatalf("Failed to create database directory %s: %v", dir, err)
			}
		}
		return sqlite.Open(sqlitePath)
	default:
		log.Fatalf("Unsupported database driver: %s (支持 sqlite/mysql/postgres)", driver)
		return nil
	}
}
