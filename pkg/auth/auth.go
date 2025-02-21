// Package auth 授权相关逻辑
package auth

import (
	"errors"
	"gohub/app/models/user"
	"gohub/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Attempt 尝试登录
func Attempt(loginID string, password string) (user.User, error) {

	// 通过 loginID 获取用户
	userModel := user.GetByMulti(loginID)
	if userModel.ID == 0 {
		return user.User{}, errors.New("用户不存在")
	}

	// 验证密码
	if !userModel.ComparePassword(password) {
		return user.User{}, errors.New("密码错误")
	}

	return userModel, nil
}

// LoginByPhone 登录指定用户
func LoginByPhone(phone string) (user.User, error) {
	userModel := user.GetByMulti(phone)
	if userModel.ID == 0 {
		return user.User{}, errors.New("用户不存在")
	}

	return userModel, nil
}

// CurrentUser 获取当前用户
func CurrentUser(c *gin.Context) user.User {
	userModel, ok := c.MustGet("current_user").(user.User)

	if !ok {
		logger.LogIf(errors.New("获取当前用户失败"))
		return user.User{}
	}

	return userModel
}

// CurrentUserID 获取当前用户 ID
func CurrentUserID(c *gin.Context) string {
	userID, ok := c.MustGet("current_user_id").(string)

	if !ok {
		logger.LogIf(errors.New("获取当前用户 ID 失败"))
		return ""
	}

	return userID
}