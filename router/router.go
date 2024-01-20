package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/synergy_api_server/controller/auth"
	"github.com/kazukiyo17/synergy_api_server/controller/game/scene"
	"github.com/kazukiyo17/synergy_api_server/middleware/jwt"
	"net/http"
)

// SetupRouter 路由信息
func SetupRouter() http.Handler {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// CORS
	router.Use(Cors())

	router.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	userApi := router.Group("/user")
	userApi.POST("/login", auth.Login)
	userApi.POST("/logout", auth.Logout)
	userApi.POST("/signup", auth.Signup)

	apiv1 := router.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/game/startScene", scene.StartScene)
		apiv1.GET("/game/scene", scene.SceneCheck)
	}

	return router
}

// Cors 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin) // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.JSON(200, gin.H{"message": "success"})
		}
		c.Next()
	}
}
