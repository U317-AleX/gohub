package config

import "gohub/pkg/config"

func init() {
	config.Add("jwt", func() map[string]interface{} {
		return map[string]interface{}{


			// 过期时间，单位是分钟，一般设置为 120 分钟
			"expire_time": config.Env("JWT_EXPIRE_TIME", 120),


			// 允许刷新时间，单位是分钟，86400 分钟是 60 天，从 token 签发开始计算
			"max_refresh_time": config.Env("JWT_MAX_REFRESH_TIME", 86400),

			// debug 模式下的过期时间，单位是分钟，一般设置为 120 分钟
			"debug_expire_time": 86400,
		}
	})
}