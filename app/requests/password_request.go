package requests

import (
	"gohub/app/requests/validators"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type ResetByPhonePasswordRequest struct {
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
		"phone":      []string{"required:手机号不能为空", "digits:手机号格式不正确"},
		"verify_code": []string{"required:验证码不能为空", "digits:验证码格式不正确"},
		"password":   []string{"required:密码不能为空", "min:密码长度不能小于 6 位"},
	}

	errs := validate(data, rules, message)

	// 检查验证码
	_data := data.(*ResetByPhonePasswordRequest)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}

// 通过邮箱重设密码
type ResetByEmailPasswordRequest struct{
	Email      string `json:"email,omitempty" valid:"email"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
	Password   string `json:"password" valid:"password"`
}

// ResetByEmail 验证表单，返回长度为0则代表验证通过
func ResetByEmail(data interface{}, c *gin.Context) map[string][]string {

	rules := govalidator.MapData{
		"email":[]string{"required", "min:4", "max:30", "email"},
		"verify_code":[]string{"required", "digits:6"},
		"password":[]string{"required", "min:6"},
	}

	message := govalidator.MapData{
		"email": []string{
            "required:Email 为必填项",
            "min:Email 长度需大于 4",
            "max:Email 长度需小于 30",
            "email:Email 格式不正确，请提供有效的邮箱地址",
        },

		"verify_code": []string{
			"required:验证码不能为空", 
			"digits:验证码格式不正确",
		},

		"password": []string{
			"required:密码不能为空", 
			"digits:密码长度不能小于 6 位",
		},
	}

	errs := validate(data, rules, message)
	
	// 检查验证码
	_data := data.(*ResetByEmailPasswordRequest)
	errs = validators.ValidateVerifyCode(_data.Email, _data.Password, errs)

	return errs
}