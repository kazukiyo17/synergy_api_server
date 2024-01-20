package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kazukiyo17/synergy_api_server/redis"
	"github.com/kazukiyo17/synergy_api_server/setting"
	"time"
)

var jwtSecret []byte

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	// 1 week 过期
	expireTime := nowTime.Add(time.Duration(setting.ServerSetting.AuthExpire) * time.Hour * 24)
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			//Issuer:    "gin-blog",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	// ParseWithClaims: parse token with claims
	// jwt.Parse: parse token without claims
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// jwt.SigningMethodHS256: signing method
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		// Valid: check if the token is valid
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func RemoveToken(token string) (err error) {
	res, err := redis.Delete(token)
	if err != nil {
		panic(err)
	}
	if !res {
		panic("token not exist")
	}
	return err
}
