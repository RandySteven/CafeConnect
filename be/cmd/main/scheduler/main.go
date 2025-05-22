package main

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apps"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/joho/godotenv"
	"log"
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

	app.PrepareJobScheduler(ctx)
}
