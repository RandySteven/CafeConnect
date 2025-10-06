package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RandySteven/CafeConnect/be/apps"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/enums"
	consumers2 "github.com/RandySteven/CafeConnect/be/handlers/consumers"
	"github.com/joho/godotenv"
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

	consumers := app.PrepareConsumer(ctx)
	runners := consumers2.InitRunner(app.Nsq)
	runners.RegisterConsumer(
		enums.TransactionTopic, consumers.TransactionConsumer.MidtransTransactionRecord,
	)
	runners.RegisterConsumer(
		enums.PaymentMidtransTopic, consumers.TransactionConsumer.MidtransPaymentConfirmation,
	)
	runners.RegisterConsumer(
		enums.OnboardingTopic, consumers.OnboardingConsumer.VerifyOnboardingToken,
	)
	runners.RegisterConsumer(
		enums.UserPointTopic, consumers.OnboardingConsumer.UserPointUpdate,
	)
	runners.RegisterConsumer(
		enums.CafeTopic, consumers.CafeConsumer.GetCafesByRadius,
	)

	_ = runners.Run(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Service is running. Press Ctrl+C to stop.")
	<-sigChan

	log.Println("Received shutdown signal, shutting down...")
	cancel()

	log.Println("Service stopped.")
}
