package auth_service

import (
	"github.com/kazukiyo17/synergy_api_server/model/auth"
	"github.com/kazukiyo17/synergy_api_server/redis"
	"github.com/kazukiyo17/synergy_api_server/utils/jwt"
)

type Auth struct {
	Username string
	Password string
	Token    string
}

// Check 检查用户登录信息
func (a *Auth) Check() (bool, error) {
	// 从redis中获取token
	rKey := "token:" + a.Token
	if redis.Exists(rKey) {
		// 已登陆
		return true, nil
	}
	// 未登陆
	return false, nil
}

func (a *Auth) Login() (bool, error) {
	// 检查Username Password是否正确
	checkSuccess, err := auth.CheckAuth(a.Username, a.Password)
	return checkSuccess, err
}

// IsUsernameExist 检查用户名是否存在
func (a *Auth) IsUsernameExist() (bool, error) {
	rkey := "username:" + a.Username
	if redis.Exists(rkey) {
		return true, nil
	}
	exist, err := auth.CheckUsername(a.Username)
	if err != nil {
		return true, err
	}
	return exist, nil
}

// Signup 注册用户
func (a *Auth) Signup() error {
	// 将用户信息写入数据库
	err := auth.AddAuth(a.Username, a.Password)
	if err != nil {
		return err
	}
	// 将用户名写入redis
	rKey := "username:" + a.Username
	err = redis.Set(rKey, a.Username, 3)
	return nil
}

// Logout 删除用户
func (a *Auth) Logout() error {
	// 删除redis中的token
	err := jwt.RemoveToken(a.Token)
	if err != nil {
		return err
	}
	return nil
}
