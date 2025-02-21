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

// 使用密码登录
type LoginByPasswordRequest struct {
	LoginID 		string `json:"login_id" valid:"login_id"`
	Password 		string `json:"password" valid:"password"`

	CaptchaID		string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer 	string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
}

// LoginByPassword 登录验证规则
func LoginByPassword(data interface{}, c *gin.Context) map[string][]string {

	rules := govalidator.MapData{
		"login_id": []string{"required", "min:3"},
		"password": []string{"required", "min:6"},
		"captcha_id": []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}

	message := govalidator.MapData{
		"login_id": []string{
			"required:登录名不能为空",
			"min:登录名长度不能少于3位",
		},
		"password": []string{
			"required:密码不能为空",
			"min:密码长度不能少于6位",
		},
		"captcha_id": []string{
			"required:图片验证码 ID 不能为空",
		},
		"captcha_answer": []string{
			"required:图片验证码答案不能为空",
			"digits:图片验证码答案长度必须为 6 位",
		},
	}

	errs := validate(data, rules, message)

	// 图片验证码
	_data := data.(*LoginByPasswordRequest)
	errs = validators.ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}