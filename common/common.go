package common

import (
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/fake_buddha_server/common/errcode"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func SendParamError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": errcode.ERROR_CODE_PARAMS_ERROR,
		"msg":  "传递参数错误：" + msg,
	})
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Resp{
		Code: errCode,
		Msg:  errcode.GetMSG(errCode),
		Data: data,
	})
	return
}
