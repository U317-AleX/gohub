package auth

import (
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/requests"
	"gohub/pkg/auth"
	"gohub/pkg/jwt"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

// LoginController 登录控制器
type LoginController struct {
	v1.BaseAPIController
}

// LoginByPhone 使用手机号和验证码登录
func (lc *LoginController) LoginByPhone(c *gin.Context) {
	// 验证表单
	request := requests.LoginByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPhone); !ok {
		return
	}

	// 验证成功, 登录
	userModel, err := auth.LoginByPhone(request.Phone)
	if err != nil {
		response.Error(c, err, "用户不存在")
		return
	}

	// 生成 token
	token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Name)

	// 返回 token
	response.JSON(c, gin.H{
		"token": token,
	})
}

// LoginByPassword 使用密码登录
func (lc *LoginController) LoginByPassword(c *gin.Context) {
	// 验证表单
	request := requests.LoginByPasswordRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPassword); !ok {
		return
	}

	// 验证成功, 登录
	userModel, err := auth.Attempt(request.LoginID, request.Password)
	if err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 生成 token
	token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Name)

	// 返回 token
	response.JSON(c, gin.H{
		"token": token,
	})
}

// RefreshToken 刷新 Access Token
func (lc *LoginController) RefreshToken(c *gin.Context) {

	// 生成新的 token
	token, err := jwt.NewJWT().RefreshToken(c)

	// 返回 token
	if err != nil {
		response.Unauthorized(c, "token 刷新失败")
	} else {
		response.JSON(c, gin.H{
			"token": token,
		})
	}
}