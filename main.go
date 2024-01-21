package main

import (
	"github.com/kazukiyo17/synergy_api_server/model"
	"github.com/kazukiyo17/synergy_api_server/redis"
	"github.com/kazukiyo17/synergy_api_server/redis_mq"
	"github.com/kazukiyo17/synergy_api_server/router"
	"github.com/kazukiyo17/synergy_api_server/setting"
	"github.com/kazukiyo17/synergy_api_server/utils/flake"
	"log"
	"net/http"
)

func init() {
	setting.Setup()
	redis_mq.Setup()
	err := redis.Setup()
	if err != nil {
		log.Printf("redis setup error: %v", err)
	}
	model.Setup()
	flake.Setup()
}

func main() {
	//conf.Load()
	routersInit := router.SetupRouter()
	server := &http.Server{
		Handler: routersInit,
		//setting.ServerSetting.HttpPort,
		Addr: ":" + setting.ServerSetting.HttpPort,
	}

	log.Printf("[info] start http server listening" + setting.ServerSetting.HttpPort)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
