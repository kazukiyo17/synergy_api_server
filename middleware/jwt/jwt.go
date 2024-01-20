package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	e "github.com/kazukiyo17/synergy_api_server/common/errcode"
	auth "github.com/kazukiyo17/synergy_api_server/service/auth"
	jwt2 "github.com/kazukiyo17/synergy_api_server/utils/jwt"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		//token := c.Query("token")
		// 从cookie中获取token
		token, err := c.Cookie("token")
		if err != nil {
			code = e.AUTH_CHECK_ERROR
		}
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			_, err := jwt2.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.AUTH_EXPIRED
				default:
					code = e.AUTH_CHECK_ERROR
				}
			} else {
				authService := auth.Auth{Token: token}
				isLogin := authService.IsLogin()
				if !isLogin {
					code = e.AUTH_EXPIRED
				}
				//rKey := "token:" + token
				//// Check if token exists in Redis
				//exists := redis.Exists(rKey)
				//if !exists {
				//	code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				//}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMSG(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
