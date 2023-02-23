package main

import (
	"bot/pkg/bot"
	"log"

	"github.com/Sagleft/utopialib-go/v2/internal/utopia"
)

func main() {
	_, err := bot.NewUtopiaBot(utopia.Config{})
	if err != nil {
		log.Fatalln(err)
	}
}
