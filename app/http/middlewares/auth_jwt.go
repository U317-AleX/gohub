// Packge : middlewares
// Auth middleware for jwt token

package middlewares

import (
	"fmt"
	"gohub/app/models/user"
	"gohub/pkg/config"
	"gohub/pkg/jwt"

	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)


func AuthJWT() gin.HandlerFunc{
	return func(c *gin.Context){
		// Get token string from request header
		tokenString, parseErr := jwt.NewJWT().ParserToken(c)

		// If failed to get token string, return error
		if parseErr != nil {
			response.Unauthorized(c, fmt.Sprintf("获取 Token 失败: %v, 请查看 %v 相关的接口认证文档", 
			parseErr, config.GetString("app.name")))
			return
		}

		// jwt 解析成功，设置用户信息
		userModel := user.Get(tokenString.UserID)

		// 如果用户不存在，返回错误
		if userModel.ID == 0 {
			response.Unauthorized(c, "用户不存在")
			return
		}

		// 将用户信息设置到请求上下文中
		// 后续的控制器可以通过 gin.context 获取用户信息
		// 这样 gin.context 就成为了一个全局的用户信息存储器
		// 就像 java 的 localThread 一样

		c.Set("current_user_id", userModel.GetStringID())
		c.Set("current_user_name", userModel.Name)
		c.Set("current_user", userModel)

		// 继续执行请求
		// 如果不主动调用 c.Next()，则请求不会继续执行
		// 也就是说，如果不调用 c.Next()，则请求到这里就结束了
		c.Next()
	}
}