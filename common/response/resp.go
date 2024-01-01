package response

import (
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/synergy_api_server/common/errcode"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  errcode.GetMSG(errCode),
		Data: data,
	})
	return
}
