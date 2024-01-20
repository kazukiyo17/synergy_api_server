package auth

import (
	e "github.com/kazukiyo17/synergy_api_server/common/errcode"
	authModel "github.com/kazukiyo17/synergy_api_server/model/auth"
	"github.com/kazukiyo17/synergy_api_server/redis"
	"github.com/kazukiyo17/synergy_api_server/service/scene"
	"github.com/kazukiyo17/synergy_api_server/setting"
	"github.com/kazukiyo17/synergy_api_server/utils/jwt"
	"log"
)

type Auth struct {
	Username string `valid:"Required; Alpha; MaxSize(20); MinSize(6);"`
	Password string `valid:"Required; MaxSize(50); MinSize(8); Match(/[A-Za-z0-9]+/)"`
	Token    string
}

// IsLogin 检查用户token
func (a *Auth) IsLogin() bool {
	rKey := "token:" + a.Token
	if redis.Exists(rKey) {
		return true
	}
	return false
}

// Login 登陆
func (a *Auth) Login() (int, string) {
	checkSuccess, err := authModel.CheckAuth(a.Username, a.Password)
	if err != nil {
		log.Printf("authService.Login() err: %v", err)
		return e.ERROR, ""
	}
	// 用户名密码错误
	if !checkSuccess {
		return e.WRONG_USERNAME_OR_PASSWORD, ""
	}
	token, err := jwt.GenerateToken(a.Username, a.Password)
	if err != nil {
		log.Printf("jwt.GenerateToken err: %v", err)
		return e.ERROR, ""
	}
	redis.Set("token:"+token, a.Username, setting.ServerSetting.AuthExpire)
	redis.Set("username:"+a.Username, "", setting.ServerSetting.AuthExpire)
	return e.SUCCESS, token
}

// IsUsernameExist 检查用户名是否存在
func (a *Auth) isUsernameExist() (bool, error) {
	rKey := "username:" + a.Username
	if redis.Exists(rKey) {
		return true, nil
	}
	exist, err := authModel.CheckUsername(a.Username)
	if err != nil {
		return true, err
	}
	return exist, nil
}

// Signup 注册用户
func (a *Auth) Signup() int {
	usernameExist, err := a.isUsernameExist()
	if err != nil {
		return e.ERROR
	}
	if usernameExist {
		return e.USERNAME_EXIST
	}
	// 注册
	err = authModel.AddAuth(a.Username, a.Password)
	if err != nil {
		return e.ERROR
	}
	// 生成初始场景
	scene.GenerateInitScene(a.Username)
	rKey := "username:" + a.Username
	redis.Set(rKey, a.Username, setting.ServerSetting.AuthExpire)
	return e.SUCCESS
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
