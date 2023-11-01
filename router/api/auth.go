package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/fake_buddha_server/common/errcode"
	"github.com/kazukiyo17/fake_buddha_server/model"
	"github.com/kazukiyo17/fake_buddha_server/utils"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// GetAuth: get auth
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := errcode.ERROR

	if ok {
		isExist := model.CheckAuth(username, password)
		if isExist {
			// GenerateToken: generate token
			token, err := utils.GenerateToken(username, password)
			if err != nil {
				code = errcode.ERROR
			} else {
				data["token"] = token
				code = errcode.SUCCESS
			}
		} else {
			code = errcode.ERROR
		}
		//} else {
		//	// MarkErrors: mark errors
		//	for _, err := range valid.Errors {
		//		//logging.Info(err.Key, err.Message)
		//	}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errcode.GetMSG(code),
		"data": data,
	})
}
