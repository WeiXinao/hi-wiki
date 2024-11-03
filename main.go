package main

import (
	"context"
	"github.com/WeiXinao/hi-wiki/ioc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func main() {
	app := ioc.InitApp()

	initPrometheus()

	otCancel := ioc.InitOTEL()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		otCancel(ctx)
	}()

	server := app.Server
	server.Run(":8080")
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}
