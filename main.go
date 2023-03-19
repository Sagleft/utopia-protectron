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
	appName        = "protectron"
	donateAddress  = "F50AF5410B1F3F4297043F0E046F205BCBAA76BEC70E936EB0F3AB94BF316804"
)

func main() {
	swissknife.PrintIntroMessage(appName, donateAddress, "CRP")

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
