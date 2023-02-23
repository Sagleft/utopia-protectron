package bot

import (
	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
)

type uBot struct {
	handler *uchatbot.ChatBot
}

type uBotConfig struct {
	WelcomeMessage string          `json:"welcomeMessage"`
	UtopiaConfig   utopiago.Config `json:"utopia"`
}

func NewUtopiaBot(cfg uBotConfig) (Bot, error) {
	b := &uBot{}

	var err error
	b.handler, err = uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: cfg.UtopiaConfig,
		Chats: []uchatbot.Chat{
			{ID: "TODO"}, // TODO: load from func args
		},
		Callbacks: uchatbot.ChatBotCallbacks{
			OnContactMessage:        b.onContactMessage,
			OnChannelMessage:        b.onChannelMessage,
			OnPrivateChannelMessage: b.onPrivateChannelMessage,

			WelcomeMessage: func(userPubkey string) string {
				return cfg.WelcomeMessage
			},
		},
		UseErrorCallback: true,
		ErrorCallback:    b.onError,
	})
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *uBot) onError(err error) {
	color.Red(err.Error())
}

func (b *uBot) onContactMessage(message structs.InstantMessage) {
	// TODO
}

func (b *uBot) onChannelMessage(message structs.WsChannelMessage) {
	// TODO
}

func (b *uBot) onPrivateChannelMessage(message structs.WsChannelMessage) {
	// TODO
}
