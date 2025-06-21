package main

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apps"
	"github.com/RandySteven/CafeConnect/be/caches"
	"github.com/RandySteven/CafeConnect/be/configs"
	consumers2 "github.com/RandySteven/CafeConnect/be/handlers/consumers"
	repositories2 "github.com/RandySteven/CafeConnect/be/repositories"
	"github.com/RandySteven/CafeConnect/be/topics"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := godotenv.Load("./files/env/.env"); err != nil {
		log.Fatalln("failed to load .env:", err)
	}
}

func main() {
	configPath, err := configs.ParseFlags()
	if err != nil {
		log.Fatalln(err)
	}

	config, err := configs.NewConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := apps.NewApps(config)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	repo := repositories2.NewRepositories(app.MySQL.Client())
	cache := caches.NewCaches(app.Redis.Client())
	topic := topics.NewTopics(app.Pub, app.Sub)

	consumers := consumers2.NewConsumers(repo, cache, topic, app.Midtrans, app.Email)

	runners := consumers2.RegisterConsumer(
		consumers.TransactionConsumer.MidtransTransactionRecord,
		consumers.OnboardingConsumer.VerifyOnboardingToken,
	)

	runners.Run(ctx)
	//go consumers.DummyConsumer.CheckHealth(ctx)
	//go consumers.TransactionConsumer.MidtransTransactionRecord(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Service is running. Press Ctrl+C to stop.")
	<-sigChan

	log.Println("Received shutdown signal, shutting down...")
	cancel()

	log.Println("Service stopped.")
}
