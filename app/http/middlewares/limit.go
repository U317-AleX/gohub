package middlewares

import (
	"gohub/pkg/app"
	"gohub/pkg/limiter"
	"gohub/pkg/logger"
	"gohub/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// LimitIP 是全局限流中间件，针对 IP 进行限流
func LimitIP(limit string) gin.HandlerFunc {
	if app.IsTesting() {
		limit = "1000000-H"
	}

	return func(c *gin.Context) {

		// 针对 IP 限流
		key := limiter.GetKeyIP(c)
		if ok := limitHandler(c, key, limit); !ok {
			return
		}

		c.Next()
	}
}

// LimitPerRoute 限流中间件，用在单独的路由中
func LimitPerRoute(limit string) gin.HandlerFunc {
	if app.IsTesting() {
		limit = "1000000-H"
	}

	return func(c *gin.Context) {

		// 针对单个路由，增加访问次数
		c.Set("limiter-once", false)

		// 针对 IP 限流
		key := limiter.GetKeyRouteWithIP(c)
		if ok := limitHandler(c, key, limit); !ok {
			return
		}

		c.Next()
	}
}

func limitHandler(c *gin.Context, key string, limit string) bool {

	// 获取超额的情况
	rate, err := limiter.CheckRate(c, key, limit)

	if err != nil {
		logger.LogIf(err)
		response.Abort500(c)
		return false
	}

	// 设置 response 的标头信息
	c.Header("X-RateLimit-Limit", cast.ToString(rate.Limit))
	c.Header("X-RateLimit-Remaining", cast.ToString(rate.Remaining))
	c.Header("X-RateLimit-Reset", cast.ToString(rate.Reset))

	if rate.Reached {

		// 提示用户超额了
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message" : "请求接口太频繁",
		})
		return false
	}

	return true
}

