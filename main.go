package main

import (
	"bot/pkg/bot"
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
)

const (
	configFilePath = "config.json"
)

func main() {
	var cfg bot.UBotConfig
	if err := swissknife.ParseStructFromJSONFile(configFilePath, &cfg); err != nil {
		log.Fatalln(err)
	}

	_, err := bot.NewUtopiaBot(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}
