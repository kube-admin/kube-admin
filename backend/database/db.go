package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 使用内存数据库，避免CGO依赖问题
	var err error
	DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库模型
	err = DB.AutoMigrate(
		&model.User{},
		&model.Cluster{},
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
