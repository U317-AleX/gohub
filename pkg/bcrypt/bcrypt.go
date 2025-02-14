//用户密码加密和验证
package bcrypt

import "golang.org/x/crypto/bcrypt"

// 生成 bcrypt 哈希
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}


// 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}