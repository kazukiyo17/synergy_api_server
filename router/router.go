package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/fake_buddha_server/controller/game/scene"
	"github.com/kazukiyo17/fake_buddha_server/middleware/jwt"
	"github.com/kazukiyo17/fake_buddha_server/router/api"
	"net/http"
)

// SetupRouter 路由信息
func SetupRouter() http.Handler {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/auth", api.GetAuth)
	router.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	apiv1 := router.Group("/api/v1")

	//apiv1.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/scene", scene.SceneInfo)
	}

	return router
}
