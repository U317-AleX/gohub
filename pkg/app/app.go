// Package app 应用信息
package app

import (
	"gohub/pkg/config"
	"time"
)

func IsLocal() bool {
    return config.Get("app.env") == "local"
}

func IsProduction() bool {
    return config.Get("app.env") == "production"
}

func IsTesting() bool {
    return config.Get("app.env") == "testing"
}

// TimeInTimeZone 获取当前时区的时间
func TimeInTimeZone() time.Time {
    loc, _ := time.LoadLocation(config.Get("app.timezone"))
    return time.Now().In(loc)
}