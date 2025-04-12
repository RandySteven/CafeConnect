package main

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/configs"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"log"
)

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
	ctx := context.Background()

	mysql, err := mysql_client.NewMySQLClient(config)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer mysql.Close()

	err = mysql.Migration(ctx)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("SUCCESS RUN MIGRATION")
}
