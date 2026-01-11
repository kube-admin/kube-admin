package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
)

// UserAPI 用户API接口
type UserAPI struct {
	userService service.UserService
}

// NewUserAPI 创建用户API实例
func NewUserAPI(userService service.UserService) *UserAPI {
	return &UserAPI{
		userService: userService,
	}
}

// ListUsers 获取用户列表
func (api *UserAPI) ListUsers(c *gin.Context) {
	users, err := api.userService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "获取用户列表失败"))
		return
	}

	// 移除密码字段后返回
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, model.SuccessResponse(users))
}

// GetUser 根据ID获取用户
func (api *UserAPI) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "无效的用户ID"))
		return
	}

	user, err := api.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(404, "用户不存在"))
		return
	}

	// 移除密码字段后返回
	user.Password = ""

	c.JSON(http.StatusOK, model.SuccessResponse(user))
}

// CreateUser 创建用户
func (api *UserAPI) CreateUser(c *gin.Context) {
	var updateReq model.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "无效的请求数据"))
		return
	}

	// 创建用户对象
	user := model.User{
		Username: updateReq.Username,
		Email:    updateReq.Email,
		Role:     updateReq.Role,
		Password: updateReq.Password, // 密码会在service层处理
	}

	if err := api.userService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(nil))
}

// UpdateUser 更新用户
func (api *UserAPI) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "无效的用户ID"))
		return
	}

	var updateReq model.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "无效的请求数据"))
		return
	}

	// 创建用户对象，只设置需要更新的字段
	user := model.User{
		ID:       uint(id),
		Username: updateReq.Username,
		Email:    updateReq.Email,
		Role:     updateReq.Role,
		Password: updateReq.Password, // 密码会在service层处理
	}

	if err := api.userService.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeleteUser 删除用户
func (api *UserAPI) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "无效的用户ID"))
		return
	}

	if err := api.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "删除用户失败"))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}
