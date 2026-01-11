package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// TestUser 测试用户模型
type TestUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func main() {
	// 测试直接使用 database/sql 打开连接
	// db, err := sql.Open("sqlite", "test.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// 测试使用 GORM 打开连接
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移模型
	err = db.AutoMigrate(&TestUser{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 创建测试用户
	testUser := TestUser{Name: "test"}
	result := db.Create(&testUser)
	if result.Error != nil {
		log.Fatal("Failed to create user:", result.Error)
	}

	fmt.Println("Success!")
}