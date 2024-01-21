package scene

import (
	"github.com/gin-gonic/gin"
	e "github.com/kazukiyo17/synergy_api_server/common/errcode"
	"github.com/kazukiyo17/synergy_api_server/common/response"
	"github.com/kazukiyo17/synergy_api_server/service/scene"
	"github.com/kazukiyo17/synergy_api_server/utils/jwt"
	"github.com/spf13/cast"
	"log"
	"net/http"
)

func SceneCheck(c *gin.Context) {
	appG := response.Gin{C: c}
	username, err := getUsernameFromToken(c)
	if err != nil {
		appG.Response(http.StatusOK, e.AUTH_CHECK_ERROR, nil)
		return
	}
	sceneId := c.Query("sceneId")
	_, err = cast.ToInt64E(sceneId)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	// 生成
	code, sceneInfo := scene.Check(sceneId, username)
	appG.Response(http.StatusOK, code, sceneInfo)
}

//func GetStartScene(c *gin.Context) {
//	appG := response.Gin{C: c}
//	token, err := c.Cookie("token")
//	if err != nil {
//		log.Printf("get token error: %v", err)
//		appG.Response(http.StatusOK, e.AUTH_CHECK_ERROR, nil)
//		return
//	}
//	claims, err := jwt.ParseToken(token)
//	if claims == nil || err != nil {
//		log.Printf("parse token error: %v", err)
//		appG.Response(http.StatusOK, e.AUTH_CHECK_ERROR, nil)
//		return
//	}
//	username := claims.Username
//	e, sceneInfo := scene.GetStartScene(username)
//	if err != nil {
//		appG.Response(http.StatusOK, e, nil)
//		return
//	}
//	appG.Response(http.StatusOK, e, sceneInfo)
//}

// InitScene 初始化
//func InitScene(c *gin.Context) {
//	appG := response.Gin{C: c}
//	username := getUsernameFromToken(c)
//	if username == "" {
//		appG.Response(http.StatusOK, e.AUTH_CHECK_ERROR, nil)
//		return
//	}
//	scene.GenerateInitScene(username)
//	appG.Response(http.StatusOK, e.SUCCESS, nil)
//}

func getUsernameFromToken(c *gin.Context) (string, error) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Printf("get token error: %v", err)
		return "", err
	}
	claims, err := jwt.ParseToken(token)
	if claims == nil || err != nil {
		log.Printf("parse token error: %v", err)
		return "", err
	}
	return claims.Username, nil
}

func StartScene(c *gin.Context) {
	appG := response.Gin{C: c}
	username, err := getUsernameFromToken(c)
	if err != nil || username == "" {
		appG.Response(http.StatusOK, e.AUTH_CHECK_ERROR, nil)
		return
	}
	initId := c.Query("chooseId")
	if initId != "1" && initId != "2" {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	code, sceneInfo := scene.GetInitScene(username, initId)
	if err != nil {
		appG.Response(http.StatusOK, code, nil)
		return
	}
	appG.Response(http.StatusOK, code, sceneInfo)
}
