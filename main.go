package main

import (
	"github.com/kazukiyo17/synergy_api_server/model"
	"github.com/kazukiyo17/synergy_api_server/redis"
	"github.com/kazukiyo17/synergy_api_server/redis_mq"
	"github.com/kazukiyo17/synergy_api_server/router"
	"github.com/kazukiyo17/synergy_api_server/setting"
	"log"
	"net/http"
)

func init() {
	setting.Setup()
	redis_mq.Setup()
	redis.Setup()
	model.Setup()
}

func main() {
	//conf.Load()
	routersInit := router.SetupRouter()
	server := &http.Server{
		Handler: routersInit,
		Addr:    ":8080",
	}

	log.Printf("[info] start http server listening 8080")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
