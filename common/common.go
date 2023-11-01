package common

import (
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/fake_buddha_server/common/errcode"
	"net/http"
)

func SendParamError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": errcode.ERROR_CODE_PARAMS_ERROR,
		"msg":  "传递参数错误：" + msg,
	})
}
