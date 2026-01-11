package model

import "time"

// User 用户模型
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 完全忽略密码字段，不在JSON响应中出现
	Email     string    `json:"email"`
	Role      string    `json:"role"` // admin, user
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID       uint   `json:"id"`
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"omitempty,min=6"` // 接收密码字段用于更新
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      User   `json:"user"`
}
