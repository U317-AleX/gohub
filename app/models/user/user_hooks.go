package user

import (
	"gohub/pkg/bcrypt"

	"gorm.io/gorm"
)

// BeforceSave 在保存之前对密码进行加密
// 通过 gorm 的钩子函数，在保存之前对密码进行加密
func (userModel *User) BeforeSave(tx *gorm.DB) (err error) {

	if !bcrypt.IsHashed(userModel.Password) {
		userModel.Password = bcrypt.HashPassword(userModel.Password)
	}
	return
}