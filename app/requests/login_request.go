package requests

import (
	"gohub/app/requests/validators"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type LoginByPhoneRequest struct {
	Phone      string `json:"phone,omitempty" valid:"phone"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
}

// LoginByPhone 登录验证规则
func LoginByPhone(data interface{}, c *gin.Context) map[string][]string {

	// 配置验证规则
	rules := govalidator.MapData{
		"phone":      []string{"required", "digits:11"},
		"verify_code": []string{"required", "digits:6"},
	}

	messages := govalidator.MapData{

		"phone":      []string{
			"required:手机号不能为空", 
			"digits:手机号格式不正确",
		},

		"verify_code": []string{
			"required:验证码不能为空", 
			"digits:验证码格式不正确",
		},
	}

	// 开始验证
	errs := validate(data, rules, messages)

	// 手机号验证
	_data := data.(*LoginByPhoneRequest)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}