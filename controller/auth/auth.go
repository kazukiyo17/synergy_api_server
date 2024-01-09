package auth

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	e "github.com/kazukiyo17/synergy_api_server/common/errcode"
	"github.com/kazukiyo17/synergy_api_server/common/response"
	"github.com/kazukiyo17/synergy_api_server/redis"
	"github.com/kazukiyo17/synergy_api_server/service/auth_service"
	"github.com/kazukiyo17/synergy_api_server/setting"
	"github.com/kazukiyo17/synergy_api_server/utils/jwt"
	"log"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; Alpha; MaxSize(20); MinSize(6);"`
	Password string `valid:"Required; MaxSize(50); MinSize(8); Match(/[A-Za-z0-9]+/)"`
}

func Logout(c *gin.Context) {
	appG := response.Gin{C: c}
	// 从cookie中获取token
	token, err := c.Cookie("jwt")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	// 登出
	authService := auth_service.Auth{Token: token}
	err = authService.Logout()
	c.SetCookie("token", "", -1, "/", setting.ServerSetting.Domain, false, true)
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func Signup(c *gin.Context) {
	appG := response.Gin{C: c}
	valid := validation.Validation{}
	// 获取参数
	username := c.PostForm("username")
	password := c.PostForm("password")
	// 验证参数
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)
	if !ok {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	// username是否已存在
	exist, err := authService.IsUsernameExist()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if exist {
		appG.Response(http.StatusUnauthorized, e.ERROR_USERNAME_EXIST, nil)
		return
	}
	// 注册
	err = authService.Signup()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func Login(c *gin.Context) {
	appG := response.Gin{C: c}
	// 获取参数
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 登录
	authService := auth_service.Auth{Username: username, Password: password}
	checkSuccess, err := authService.Login()
	if err != nil {
		log.Printf("authService.Login() err: %v", err)
		appG.Response(http.StatusOK, e.ERROR_AUTH, nil)
		return
	}
	// 用户名密码错误
	if !checkSuccess {
		appG.Response(http.StatusOK, e.WRONG_USERNAME_OR_PASSWORD, nil)
		return
	}
	// token 写入cookie, 3天过期
	token, err := jwt.GenerateToken(authService.Username, authService.Password)
	if err != nil {
		log.Printf("jwt.GenerateToken err: %v", err)
		appG.Response(http.StatusOK, e.ERROR_AUTH, nil)
		return
	}
	rKey := "token:" + token
	err = redis.Set(rKey, username, 3)
	if err != nil {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
	c.SetCookie("token", token, 3*24*3600, "/", setting.ServerSetting.Domain, false, true)
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
