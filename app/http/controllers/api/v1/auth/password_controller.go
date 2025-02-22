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

// ResetPassword 重置密码
func (p *PasswordController) ResetPassword(c *gin.Context) {
	// 验证表单
	request := requests.ResetPasswordRequest{}
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