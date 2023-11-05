package user

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/fake_buddha_server/common"
	e "github.com/kazukiyo17/fake_buddha_server/common/errcode"
	"github.com/kazukiyo17/fake_buddha_server/utils/jwt"
	"net/http"
)

type Info struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	appG := common.Gin{C: c}
	valid := validation.Validation{}

	username := c.PostForm("username")
	password := c.PostForm("password")

	a := Info{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		//app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	//authService := auth_service.Auth{Token: token}
	//isExist, err := authService.Check()
	//if err != nil {
	//	appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
	//	return
	//}

	//if !isExist {
	//	appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
	//	return
	//}

	token, err := jwt.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
