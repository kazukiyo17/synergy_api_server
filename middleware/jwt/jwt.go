package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	e "github.com/kazukiyo17/synergy_api_server/common/errcode"
	auth "github.com/kazukiyo17/synergy_api_server/service/auth"
	jwt2 "github.com/kazukiyo17/synergy_api_server/utils/jwt"
	"log"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		code = e.SUCCESS
		token, err := c.Cookie("token")

		if err != nil {
			log.Printf("get token error: %v", err)
			code = e.AUTH_CHECK_ERROR
		}
		if token == "" {
			log.Printf("token is empty")
			code = e.AUTH_CHECK_ERROR
		} else {
			_, err = jwt2.ParseToken(token)
			if err != nil {
				log.Printf("parse token error: %v", err)
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
