package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接。dbPath 为 SQLite 文件路径，空则使用内存数据库。
// 采用文件持久化，确保用户与集群配置在重启后不丢失。
func InitDB(dbPath string) {
	if dbPath == "" {
		dbPath = "file::memory:?cache=shared"
	} else {
		// 确保数据目录存在
		if dir := filepath.Dir(dbPath); dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				log.Fatalf("Failed to create database directory %s: %v", dir, err)
			}
		}
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库模型
	err = DB.AutoMigrate(
		&model.User{},
		&model.Cluster{},
		&model.AuditLog{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 检查是否有管理员用户，如果没有则创建一个
	var adminCount int64
	DB.Model(&model.User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		// 创建默认管理员用户
		adminUser := model.User{
			Username: "admin",
			Email:    "admin@example.com",
			Role:     "admin",
		}
		// 使用bcrypt加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}
		adminUser.Password = string(hashedPassword)

		// 保存到数据库
		if err := DB.Create(&adminUser).Error; err != nil {
			log.Fatal("Failed to create admin user:", err)
		}
		log.Println("Created default admin user: admin/admin123")
	}

	log.Println("Database initialized successfully")
}
