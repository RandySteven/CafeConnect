package main

import (
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apps"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/enums"
	"github.com/joho/godotenv"
	"log"
	"os/exec"
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

	app, err := apps.NewApps(config)
	if err != nil {
		log.Fatalln(`Error starting app `, err)
		return
	}

	cmd := ``

	for cmd != `EXIT` {
		fmt.Print(">>Action [CREATE | DROP | SEE] : ")
		fmt.Scanln(&cmd)

		switch cmd {
		case `CREATE`:
			prepareTopics(app)
		case `DROP`:
			cmdDrop := `
				for topic in $(kafka-topics --bootstrap-server localhost:9092 --list); do 
				  kafka-topics --bootstrap-server localhost:9092 --delete --topic "$topic"
				done
			`
			command := exec.Command("bash", "-c", cmdDrop)
			output, err := command.CombinedOutput()
			if err != nil {
				log.Fatalf("Error deleting topics: %v\nOutput: %s", err, string(output))
			}
		case `SEE`:
			seeTopic(app)
		}
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
