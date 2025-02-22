// Package routes 注册路由
package routes

import (
	"gohub/app/http/controllers/api/v1/auth"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {

    // 测试一个 v1 的路由组，我们所有的 v1 版本的路由都将存放到这里
    v1 := r.Group("/v1")
    {
        authGroup := v1.Group("/auth")
        {
            suc := new(auth.SignupController)

            // 判断手机是否已注册
            authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)

            // 判断 Email 是否已注册
            authGroup.POST("/signup/email/exist", suc.IsEmailExist)
            
            //验证码Controller对象
            vcc := new(auth.VerifyCodeController)
            
            //注册获取验证码图片的路径
            authGroup.POST("/verify-codes/captcha", vcc.ShowCaptcha)

            authGroup.POST("/verify-codes/phone", vcc.SendUsingPhone)
            
            authGroup.POST("/verify-codes/email", vcc.SendUsingEmail)

            authGroup.POST("/signup/using-phone", suc.SignupUsingPhone)

            authGroup.POST("/signup/using-email", suc.SignupUsingEmail)

            // 登录
            lgc := new(auth.LoginController)

            // 使用手机号和验证码登录
            authGroup.POST("/login/using-phone", lgc.LoginByPhone)

            // 使用密码登录
            authGroup.POST("/login/using-password", lgc.LoginByPassword)

            // 刷新 token
            authGroup.POST("/login/refresh-token", lgc.RefreshToken)

            pwc := new(auth.PasswordController)

            // 重置密码
            authGroup.POST("/password-reset/using-phone", pwc.ResetPassword)
        }
    }
}