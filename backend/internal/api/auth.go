package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/database"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
	"golang.org/x/crypto/bcrypt"
)

// AuthAPI 认证API
type AuthAPI struct {
	userService service.UserService
}

// NewAuthAPI 创建认证API实例
func NewAuthAPI(userService service.UserService) *AuthAPI {
	return &AuthAPI{
		userService: userService,
	}
}

// Login 登录
func (a *AuthAPI) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "请求参数错误"))
		return
	}

	// 查询数据库中的用户
	var user model.User
	// 使用GORM的结构体查询方式，避免手动编写SQL字段名
	result := database.DB.Where(&model.User{Username: req.Username}).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse(401, "用户名不存在"))
		return
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse(401, "密码错误"))
		return
	}

	// 生成Token
	token, expiresAt, err := model.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "生成Token失败"))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}))
}

// GetUserInfo 获取用户信息
func (a *AuthAPI) GetUserInfo(c *gin.Context) {
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	user := model.User{
		Username: username.(string),
		Role:     role.(string),
	}

	c.JSON(http.StatusOK, model.SuccessResponse(user))
}
