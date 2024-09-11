package api

import (
	"github.com/Ntrashh/crawlerctl/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// LoginHandler 处理登录请求并生成 JWT Token
func LoginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 绑定并验证请求体
	if err := c.ShouldBindJSON(&loginData); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 从配置中获取用户名和密码
	expectedUsername := config.AppConfig.Auth.Username
	expectedPassword := config.AppConfig.Auth.Password

	// 校验用户名和密码
	if loginData.Username == expectedUsername && loginData.Password == expectedPassword {
		// 生成 JWT Token
		expirationTime := time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireHours) * time.Hour)
		claims := &JWTClaims{
			Username: loginData.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "生成 Token 失败")
			return
		}
		SuccessResponse(c, gin.H{"token": tokenString})
	} else {
		ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
	}
}
