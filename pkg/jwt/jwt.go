// Package jwt 处理 JWT 认证
package jwt

import (
	"errors"
	"strings"
	"time"

	"gohub/pkg/config"
	"gohub/pkg/app"

	"gohub/pkg/logger"
	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt"
)

var (
	ErrTokenExpired 			error = errors.New("token is expired")
	ErrTokenExpiredMaxRefresh 	error = errors.New("token is expired, max refresh reached")
	ErrTokenMalformed 			error = errors.New("that's not even a token")
	ErrTokenInvalid 			error = errors.New("token is invalid")
	ErrHeaderEmpty 				error = errors.New("header is empty")
	ErrHeaderMalformed 			error = errors.New("header is malformed")
)

// JWT 定义一个 JWT 结构体
type JWT struct {

	// SigningKey 签名密钥
	SigningKey []byte

	// MaxRefresh token 最大刷新时间
	MaxRefresh time.Duration
}

// JWTCustomClaims 自定义荷载
type JWTCustomClaims struct {
	UserID 		string `json:"user_id"`
	Username 	string `json:"username"`
	ExpiresAt 	int64 `json:"expires_at"`

	// jwtpkg.StandardClaims 包含了一些标准字段
	// 这里只是简单的嵌入，如果有需要可以自定义更多字段
	// 例如：Issuer, Subject, Audience 等

	jwtpkg.StandardClaims
}

// NewJWT 创建一个 JWT 实例
func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(config.GetString("app.key")),
		MaxRefresh: time.Duration(config.GetInt("jwt.max_refresh_time")) * time.Minute,
	}
}

// ParseToken 解析 token, 这个函数会在中间件中使用
func (jwt *JWT) ParserToken(c *gin.Context) (*JWTCustomClaims, error) {

	// 从请求头中获取 token 字符串
	tokenString, parseErr := jwt.getTokenFromHeader(c)

	// 如果获取 token 字符串失败，返回错误
	if parseErr != nil {
		return nil, parseErr
	}

	// 解析 token 字符串
	token, err := jwt.parseTokenString(tokenString)

	// 如果解析 token 字符串失败，返回错误
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		if ok {
			if validationErr.Errors == jwtpkg.ValidationErrorExpired {
				// Token 过期
				return nil, ErrTokenExpired
			} else if validationErr.Errors == jwtpkg.ValidationErrorMalformed {
				// Token 格式错误
				return nil, ErrTokenMalformed
			}
		}
		
		// Token 无效
		return nil, ErrTokenInvalid
	}

	// 解析 token 成功，返回解析结果
	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	// 否则，返回 Token 无效
	return nil, ErrTokenInvalid
}

func (jwt *JWT) RefreshToken(c *gin.Context) (string, error) {

	// 从请求头中获取 token 字符串
	tokenString, parseErr := jwt.getTokenFromHeader(c)

	// 如果获取 token 字符串失败，返回错误
	if parseErr != nil {
		return "", parseErr
	}

	// 解析 token 字符串
	token, err := jwt.parseTokenString(tokenString)

	// 如果解析 token 字符串失败，返回错误
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		// 如果是 Token 过期错误
		// 并且 token 的刷新时间小于最大刷新时间
		// 则刷新 token
		// 否则返回错误
		if !ok || validationErr.Errors != jwtpkg.ValidationErrorExpired {
			return "", err
		}
	}

	// 解析 token 荷载数据
	claims := token.Claims.(*JWTCustomClaims)

	// 判断是否超过了最大刷新时间
	// 如果超过了最大刷新时间，返回错误

	// time_in_timezone - issue_at < max_refresh
	// issue_at > time_in_timezone - max_refresh
	x := app.TimeInTimeZone().Add(-jwt.MaxRefresh).Unix()
	if claims.IssuedAt > x {
		// 修改过期时间
		claims.StandardClaims.ExpiresAt = jwt.expireAtTime()
		return jwt.createToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
}


// IssueToken 签发 token, 这个函数会在登录成功后调用
func (jwt *JWT) IssueToken(useID string, userName string) string {

	// 构造 JWT 荷载数据
	expireAtTime := jwt.expireAtTime()

	claims := JWTCustomClaims{
		UserID: useID,
		Username: userName,
		ExpiresAt: expireAtTime,

		StandardClaims: jwtpkg.StandardClaims{
			NotBefore: 	app.TimeInTimeZone().Unix(), // 生效时间
			IssuedAt:	app.TimeInTimeZone().Unix(), // 签发时间
			ExpiresAt: 	expireAtTime,				 // 过期时间
			Issuer: 	config.GetString("app.name"),// 签发者
		},
	}

	// 创建 token
	token, err := jwt.createToken(claims)
	if  err != nil {
		logger.LogIf(err)
		return ""
	}

	return token
}

// createToken 创建 token
func (jwt *JWT) createToken(claims JWTCustomClaims) (string, error) {

	// 创建一个 token 实例
	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)

	// 签名 token
	return token.SignedString(jwt.SigningKey)
}

// expireAtTime 计算过期时间
func (jwt *JWT) expireAtTime() int64 {
	
	// 获取当前时区时间
	timenow := app.TimeInTimeZone()
	
	var expireTime int64
	
	if config.GetBool("app.debug") {
		expireTime = config.GetInt64("jwt.debug_expire_time")
		} else {
			expireTime = config.GetInt64("jwt.expire_time")
		}
		
	expire := time.Duration(expireTime) * time.Minute

	// 过期时间 = 当前时间 + token 过期时间
	return timenow.Add(expire).Unix()
}

// parseTokenString 解析 token 字符串
func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	
	// 解析 token 字符串
	return jwtpkg.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.SigningKey, nil
	})
}

// getTokenFromHeader 从请求头中获取 token 字符串
// Authorization: Bearer ...
func (jwt *JWT) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}

	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if  !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrHeaderMalformed
	}

	return parts[1], nil
}

