package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kazukiyo17/synergy_api_server/redis"
	"time"
)

const (
	TOKEN_EXPIRE_TIME = 3
)

var jwtSecret []byte

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour) // 3 hours 过期
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

func AddToken(username, password string) (token string) {
	// 生成token
	token, err := GenerateToken(username, password)
	if err != nil {
		panic(err)
	}
	// 将token写入redis, 3天过期
	err = redis.Set(token, username, TOKEN_EXPIRE_TIME)
	if err != nil {
		panic(err)
	}
	return token
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
