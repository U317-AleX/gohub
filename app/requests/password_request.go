package requests

import (
	"gohub/app/requests/validators"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type ResetPasswordRequest struct {
	Phone      string `json:"phone,omitempty" valid:"phone"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
	Password   string `json:"password" valid:"password"`
}

// ResetByPhone 验证表单，返回长度为 0 则验证通过
func ResetByPhone(data interface{}, c *gin.Context) map[string][]string {

	rules := govalidator.MapData{
		"phone":      []string{"required", "digits:11"},
		"verify_code": []string{"required", "digits:6"},
		"password":   []string{"required", "min:6"},
	}

	message := govalidator.MapData{
		"phone":      []string{"手机号不能为空", "手机号格式不正确"},
		"verify_code": []string{"验证码不能为空", "验证码格式不正确"},
		"password":   []string{"密码不能为空", "密码长度不能小于 6 位"},
	}

	errs := validate(data, rules, message)

	// 检查验证码
	_data := data.(*ResetPasswordRequest)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}