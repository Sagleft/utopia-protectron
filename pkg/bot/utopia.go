package bot

import (
	"bot/pkg/memory"
	"fmt"

	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
)

const defaultAccountName = "account.db"

type uBot struct {
	handler *uchatbot.ChatBot
	dbConn  memory.Memory
}

type UBotConfig struct {
	WelcomeMessage string          `json:"welcomeMessage"`
	AccountName    string          `json:"accountName"`
	UtopiaConfig   utopiago.Config `json:"utopia"`
}

func NewUtopiaBot(cfg UBotConfig, db memory.Memory) (Bot, error) {
	b := &uBot{dbConn: db}

	var err error
	b.handler, err = uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: cfg.UtopiaConfig,
		Chats:  []uchatbot.Chat{},
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

	return b, b.fixAccountName(cfg.AccountName)
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

func (b *uBot) fixAccountName(accountName string) error {
	data, err := b.handler.GetOwnContact()
	if err != nil {
		return fmt.Errorf("get own contact: %w", err)
	}

	if data.Nick == defaultAccountName {
		if err := b.handler.SetAccountNickname(accountName); err != nil {
			return fmt.Errorf("set account nickname: %w", err)
		}
	}
	return nil
}
