// Package requests 处理请求数据和表单验证
package requests

import (
	"gohub/app/requests/validators"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type SignupPhoneExistRequest struct {
    Phone string `json:"phone,omitempty" valid:"phone"`
}

func SignupPhoneExist(data interface{}, c *gin.Context) map[string][]string {

    // 自定义验证规则
    rules := govalidator.MapData{
        "phone": []string{"required", "digits:11"},
    }

    // 自定义验证出错时的提示
    messages := govalidator.MapData{
        "phone": []string{
            "required:手机号为必填项，参数名称 phone",
            "digits:手机号长度必须为 11 位的数字",
        },
    }

    return validate(data, rules, messages)
}


type SignupEmailExistRequest struct {
    Email string `json:"email,omitempty" valid:"email"`
}

func SignupEmailExist(data interface{}, c *gin.Context) map[string][]string {

    // 自定义验证规则
    rules := govalidator.MapData{
        "email": []string{"required", "min:4", "max:30", "email"},
    }

    // 自定义验证出错时的提示
    messages := govalidator.MapData{
        "email": []string{
            "required:Email 为必填项",
            "min:Email 长度需大于 4",
            "max:Email 长度需小于 30",
            "email:Email 格式不正确，请提供有效的邮箱地址",
        },
    }

    return validate(data, rules, messages)
}

// SignupUsingPhoneRequest 通过手机注册的请求信息
type SignupUsingPhoneRequest struct{
    Phone           string `json:"phone,omitempty" valid:"phone"`
    Name            string `json:"name,omitempty" valid:"name"`
    Password        string `json:"password" valid:"password"`
    PasswordConfirm string `json:"password_confirm,omitempty" valid:"password_confirm"`
    VerifyCode      string `json:"verify_code,omitempty" valid:"verify_code"`
}

//使用手机登录
func SignupUsingPhone(data interface{}, c *gin.Context) map[string][]string{

    //配置格式验证规则
    rules := govalidator.MapData{
        "phone" : []string{"required", "digits:11", "not_exists:user,phone"},
        "name" : []string{"required", "alpha_num", "between:3,20", "not_exists:user,name"},
        "password" : []string{"required", "min:6"},
        "password_confirm" : []string{"required"},
        "verify_code" : []string{"required", "digits:6"},
    }

    //设置验证未通过时的错误信息
    messages := govalidator.MapData{
        "phone" : []string{
            "required : 手机号必须填写",
            "digit : 手机号长度必须为11位数字",
        },

        "name" : []string{
            "required : 用户名必须填写",
            "alpha_num : 用户名格式错误, 只允许字母和数字",
            "between: 用户名长度须在3到20位之间",
        },

        "password" : []string{
            "required : 用户密码必须填写",
            "min : 用户密码不少于六位数",
        },

        "password_confirm" : []string{
            "required : 请确认密码",
        },

        "verify_code" : []string{
            "required : 验证码必须填写",
            "digit : 验证码必须为六位数",
        },
    }

    //格式验证
    errs := validate(data, rules, messages)

    _data := data.(*SignupUsingPhoneRequest)

    //验证密码是否输入正确
    errs = validators.ValidatePasswordConfirm(_data.Password, _data.PasswordConfirm, errs)

    //验证验证码是否正确
    errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

    return errs
}