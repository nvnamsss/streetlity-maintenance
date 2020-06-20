package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"streetlity-maintenance/model"
	"streetlity-maintenance/router"
	"streetlity-maintenance/server"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var Router *mux.Router = mux.NewRouter()
var Server http.Server

func main() {
	model.Connect()
	var wait time.Duration

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	loggedRouter := handlers.LoggingHandler(os.Stdout, Router)

	router.Handle(Router)
	Server := &http.Server{
		Addr:         "0.0.0.0:9002",
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      loggedRouter,
	}

	go func() {
		if err := Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	go server.Create()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)

	defer cancel()

	Server.Shutdown(ctx)
	log.Println("shutting down")

	os.Exit(0)
}
