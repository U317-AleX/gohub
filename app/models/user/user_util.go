package user

import (
    "gohub/pkg/database"
)

// IsEmailExist 判断 Email 已被注册
func IsEmailExist(email string) bool {
    var count int64
    database.DB.Model(User{}).Where("email = ?", email).Count(&count)
    return count > 0
}

// IsPhoneExist 判断手机号已被注册
func IsPhoneExist(phone string) bool {
    var count int64
    database.DB.Model(User{}).Where("phone = ?", phone).Count(&count)
    return count > 0
}

// GetByPhone 通过手机号获取用户
func GetByPhone(phone string) User {
    var user User
    database.DB.Where("phone = ?", phone).First(&user)
    return user
}

// GetByMulti 通过手机号/邮箱/用户名获取用户
func GetByMulti(loginID string) User {
    var user User
    database.DB.Where("phone = ? OR email = ? OR name = ?", loginID, loginID, loginID).First(&user)
    return user
}