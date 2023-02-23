package main

import (
	"bot/pkg/bot"
	"bot/pkg/memory"
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
)

const (
	configFilePath = "config.json"
	dbFilename     = "memory.db"
)

func main() {
	var cfg bot.UBotConfig
	if err := swissknife.ParseStructFromJSONFile(configFilePath, &cfg); err != nil {
		log.Fatalln(err)
	}

	db, err := memory.NewLocalDB(dbFilename)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := bot.NewUtopiaBot(cfg, db); err != nil {
		log.Fatalln(err)
	}

	log.Println("bot started")
	swissknife.RunInBackground()
}
