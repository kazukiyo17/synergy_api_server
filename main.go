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
	err := redis.Setup()
	if err != nil {
		log.Fatalf("redis setup error: %v", err)
	}
	model.Setup()
}

func main() {
	//conf.Load()
	routersInit := router.SetupRouter()
	server := &http.Server{
		Handler: routersInit,
		Addr:    ":443",
	}

	log.Printf("[info] start http server listening 443")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
