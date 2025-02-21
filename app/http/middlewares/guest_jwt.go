// Guest middleware 用在一些只有游客才能访问的接口上
// 例如登录、注册等接口
package middlewares

import (
	"gohub/pkg/jwt"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

// GuestJWT 强制使用游客身份访问
func GuestJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if  len(c.GetHeader("Authorization")) > 0 {

			_, err := jwt.NewJWT().ParserToken(c)
			if err == nil {
				response.Unauthorized(c, "请使用游客身份访问")
				c.Abort()
				return
			}
		}

		// 继续执行请求
		c.Next()
	}
}