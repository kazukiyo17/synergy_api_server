package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	e "github.com/kazukiyo17/synergy_api_server/common/errcode"
	"github.com/kazukiyo17/synergy_api_server/service/auth_service"
	jwt2 "github.com/kazukiyo17/synergy_api_server/utils/jwt"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			_, err := jwt2.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			} else {
				authService := auth_service.Auth{Token: token}
				isExist, err := authService.Check()
				if err != nil {
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
				if !isExist {
					code = e.ERROR_AUTH_EXPIRED
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
