// Package auth 授权相关逻辑
package auth

import (
	"errors"
	"gohub/app/models/user"
)

// Attempt 尝试登录
func Attempt(email string, password string) (user.User, error) {

	// 通过 email 获取用户
	userModel := user.GetByMulti(email)
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