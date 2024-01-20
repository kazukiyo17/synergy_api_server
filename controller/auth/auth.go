package auth

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	e "github.com/kazukiyo17/synergy_api_server/common/errcode"
	"github.com/kazukiyo17/synergy_api_server/common/response"
	authService "github.com/kazukiyo17/synergy_api_server/service/auth"
	"github.com/kazukiyo17/synergy_api_server/setting"
	"net/http"
)

type validAuth struct {
	Username string `valid:"Required; Alpha; MaxSize(20); MinSize(6);"`
	Password string `valid:"Required; MaxSize(50); MinSize(8); Match(/[A-Za-z0-9]+/)"`
}

func Logout(c *gin.Context) {
	appG := response.Gin{C: c}
	// 从cookie中获取token
	token, err := c.Cookie("jwt")
	if err != nil {
		appG.Response(http.StatusOK, e.AUTH_CHECK_ERROR, nil)
		return
	}
	// 登出
	service := authService.Auth{Token: token}
	err = service.Logout()
	c.SetCookie("token", "", -1, "/", setting.ServerSetting.Domain, false, true)
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func Signup(c *gin.Context) {
	appG := response.Gin{C: c}
	valid := validation.Validation{}
	var service authService.Auth
	err := c.BindJSON(&service)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	ok, err := valid.Valid(&service)
	if !ok || err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	//service := authService.Auth{Username: username, Password: password}
	// 注册
	code := service.Signup()
	appG.Response(http.StatusOK, code, nil)
}

func Login(c *gin.Context) {
	appG := response.Gin{C: c}
	// 获取参数
	//username := c.PostForm("username")
	//password := c.PostForm("password")
	var service authService.Auth
	err := c.BindJSON(&service)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	// 登录
	//service := authService.Auth{Username: username, Password: password}
	res, token := service.Login()
	c.SetCookie("token", token, 3*24*3600, "/", setting.ServerSetting.Domain, false, false)
	appG.Response(http.StatusOK, res, nil)
}
