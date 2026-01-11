package service

import (
	"errors"
	"time"

	"github.com/kube-admin/kube-admin/backend/database"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	ListUsers() ([]model.User, error)
	GetUserByID(id uint) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
}

// userService 用户服务实现
type userService struct{}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{}
}

// ListUsers 获取用户列表
func (s *userService) ListUsers() ([]model.User, error) {
	var users []model.User
	result := database.DB.Find(&users)
	return users, result.Error
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(user *model.User) error {
	// 检查用户名是否已存在
	var existingUser model.User
	if result := database.DB.Where("username = ?", user.Username).First(&existingUser); result.Error == nil {
		return errors.New("用户名已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 设置创建和更新时间
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// 保存到数据库
	result := database.DB.Create(user)
	return result.Error
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(user *model.User) error {
	// 查找现有用户
	existingUser, err := s.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	// 如果密码不为空，更新密码
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		existingUser.Password = string(hashedPassword)
	}

	// 更新其他字段
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.Role = user.Role
	existingUser.UpdatedAt = time.Now()

	// 保存到数据库
	result := database.DB.Save(existingUser)
	return result.Error
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	result := database.DB.Delete(&model.User{}, id)
	return result.Error
}