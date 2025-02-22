package auth

import (
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/requests"
	"gohub/pkg/response"
	"gohub/app/models/user"

	"github.com/gin-gonic/gin"
)

type PasswordController struct {
	v1.BaseAPIController
}

// ResetByPhone 重置密码
func (p *PasswordController) ResetByPhone(c *gin.Context) {
	// 验证表单
	request := requests.ResetByPhonePasswordRequest{}
	if ok := requests.Validate(c, &request, requests.ResetByPhone); !ok {
		return
	}

	// 重置密码
	userModel := user.GetByPhone(request.Phone)
	if userModel.ID == 0 {
		response.Abort404(c)
	} else {
		userModel.Password = request.Password
		userModel.Save()
		response.Success(c)
	}
}

// ResetByEmail 重置密码
func (p *PasswordController) ResetByEmail(c *gin.Context) {
	// 验证表单
	request := requests.ResetByEmailPasswordRequest{}
	if ok := requests.Validate(c, &request, requests.ResetByEmail); !ok {
		return
	}

	// 重置密码
	userModel := user.GetByMulti(request.Email)
	if userModel.ID == 0 {
		response.Abort404(c)
	} else {
		userModel.Password = request.Password
		userModel.Save()
		response.Success(c)
	}
}