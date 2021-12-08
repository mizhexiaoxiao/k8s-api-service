package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/config"
	"github.com/mizhexiaoxiao/k8s-api-service/models"
	"github.com/mizhexiaoxiao/k8s-api-service/routers"
)

func init() {
	config.Setup()
	models.Setup()
	app.Setup()
}

func main() {
	routersInit := routers.InitRouter()
	server := &http.Server{
		Addr:         config.AppAddr(),
		Handler:      routersInit,
		ReadTimeout:  time.Duration(config.ReadTimeout()),
		WriteTimeout: time.Duration(config.WriteTimeout()),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server err: %v", err)
	}
}
