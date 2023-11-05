package auth_service

import (
	"github.com/kazukiyo17/fake_buddha_server/controller/user"
	"github.com/kazukiyo17/fake_buddha_server/utils/redis"
)

type Auth struct {
	userId int64
	token  string
}

func (a *Auth) C  eck() bool {
	return redis.Exists(a.token)
}

func (a *Auth) Save() bool {
	err := redis.Set(a.token, a.userId, 3600)
	if err != nil {
		return false
	}
	return true
}

func (a *Auth) Login(info user.Info, token string) (string, error) {
	// 检查redis
	if redis.Exists(token) {
		// 重置过期时间
		redis.Set(token, a.userId, 3600)
		return token, nil
	}
}
