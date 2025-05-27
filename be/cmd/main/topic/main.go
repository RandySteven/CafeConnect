package main

import (
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apps"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/enums"
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
	cmd := ``
	fmt.Print(">>Action [CREATE | DROP | SEE] : ")
	fmt.Scanln(&cmd)
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

	app, err := apps.NewApps(config)
	if err != nil {
		log.Fatalln(`Error starting app `, err)
		return
	}

	switch cmd {
	case `CREATE`:
		prepareTopics(app)
	case `DROP`:
		clearAllTopics(app)
	case `SEE`:
		seeTopic(app)
	}

}

func prepareTopics(app *apps.App) {
	err := app.Kafka.RegisterTopics(
		enums.DummyTopic,
		enums.TransactionTopic,
		enums.ProductTopic,
		enums.OnboardingTopic,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func clearAllTopics(app *apps.App) {
	if err := app.Kafka.ClearAllTopics(); err != nil {
		log.Println(err)
		return
	}
}

func seeTopic(app *apps.App) {
	app.Kafka.ReadTopics()
}
