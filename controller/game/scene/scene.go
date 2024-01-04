package scene

import (
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/synergy_api_server/service/scene_service"
	"github.com/kazukiyo17/synergy_api_server/utils/jwt"
	"github.com/spf13/cast"
	"net/http"
	"strings"
)

func Check(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		// 鉴权失败
		c.JSON(http.StatusUnauthorized, gin.H{"message": "auth failed"})
		return
	}
	claims, err := jwt.ParseToken(token)
	if claims == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "auth failed"})
		return
	}
	sceneUrl := c.Query("url")
	if sceneUrl == "" {
		c.JSON(-1, gin.H{"message": "url is empty"})
		return
	}
	sceneIdStr := sceneUrl[strings.LastIndex(sceneUrl, "/")+1 : strings.LastIndex(sceneUrl, ".")]
	// 如果为start.txt, 或end.txt, 直接返回
	if sceneIdStr == "start" || sceneIdStr == "end" {
		c.JSON(200, gin.H{"message": "success"})
		return
	}
	// 是否为数字
	_, err = cast.ToInt64E(sceneIdStr)
	if err != nil {
		c.JSON(-1, gin.H{"message": "sceneId is invalid"})
		return
	}
	// 用户是否有权限
	isCreator := scene_service.CheckSceneCreator(sceneIdStr, claims.Username)
	if !isCreator {
		c.JSON(-1, gin.H{"message": "permission denied"})
		return
	}
	// 生成
	err = scene_service.Check(sceneIdStr, claims.Username)
	if err != nil {
		c.JSON(-1, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
