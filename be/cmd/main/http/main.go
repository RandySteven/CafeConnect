package main

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apps"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := godotenv.Load("./files/env/.env")
	if err != nil {
		log.Fatalln(`failed to load .env `, err)
		return
	}
}

func main() {
	configPath, err := configs.ParseFlags()
	if err != nil {
		log.Fatalln(err)
		return
	}

	config, err := configs.NewConfig(configPath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	ctx := context.TODO()

	app, err := apps.NewApps(config)
	if err != nil {
		log.Fatalln(`Error starting app `, err)
		return
	}

	apis := app.PrepareHttpHandler(ctx)
	r := mux.NewRouter()
	router := routes.NewEndpointRouters(apis)
	routes.InitRouter(router, r)

	go config.Run(r)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = app.RefreshRedis(ctx); err != nil {
		log.Fatal(err)
		return
	}
}
