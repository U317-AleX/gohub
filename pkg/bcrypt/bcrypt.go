// 用户密码加密和验证
package bcrypt

import (
	"gohub/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// 生成 bcrypt 哈希
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	logger.LogIf(err)

	return string(hash)
}


// 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 判断字符串是否是哈希后的密码
func IsHashed(password string) bool {
	// bcrypt 加密后的密码长度是 60
	return len(password) == 60
}