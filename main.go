package main

import (
	"github.com/kazukiyo17/fake_buddha_server/common/conf"
	"github.com/kazukiyo17/fake_buddha_server/router"
	"github.com/kazukiyo17/fake_buddha_server/setting"
	"log"
	"net/http"
)

func init() {
	setting.Setup()
}

func main() {
	conf.Load()
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
