package middleware

import (
	"crawlerctl/api"
	"crawlerctl/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// JWTClaims 定义 JWT 的声明部分
type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			api.ErrorResponse(c, http.StatusUnauthorized, "缺少 Authorization 头")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			api.ErrorResponse(c, http.StatusUnauthorized, "无效的 Authorization 头")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 Token
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			api.ErrorResponse(c, http.StatusUnauthorized, "无效的 Token")
			c.Abort()
			return
		}
		// 设置用户名到上下文
		c.Set("username", claims.Username)
		c.Next()
	}
}
