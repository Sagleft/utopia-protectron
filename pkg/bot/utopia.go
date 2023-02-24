package bot

import (
	"bot/pkg/memory"
	"fmt"

	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
)

const (
	channelIDLength    = 32
	defaultAccountName = "account.db"
)

type uBot struct {
	handler *uchatbot.ChatBot
	dbConn  memory.Memory
	cfg     UBotConfig
}

type UBotConfig struct {
	WelcomeMessage string          `json:"welcomeMessage"`
	AccountName    string          `json:"accountName"`
	UtopiaConfig   utopiago.Config `json:"utopia"`
}

func NewUtopiaBot(cfg UBotConfig, db memory.Memory) (Bot, error) {
	b := &uBot{
		dbConn: db,
		cfg:    cfg,
	}

	var err error
	b.handler, err = uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: cfg.UtopiaConfig,
		Chats:  []uchatbot.Chat{},
		Callbacks: uchatbot.ChatBotCallbacks{
			OnContactMessage:        b.onContactMessage,
			OnChannelMessage:        b.onChannelMessage,
			OnPrivateChannelMessage: b.onPrivateChannelMessage,

			WelcomeMessage: b.onWelcomeMessage,
		},
		UseErrorCallback: true,
		ErrorCallback:    b.onError,
	})
	if err != nil {
		return nil, err
	}

	return b, b.fixAccountName(cfg.AccountName)
}

func (b *uBot) onWelcomeMessage(userPubkey string) string {
	isUserSaved, err := b.dbConn.IsUserExists(memory.User{Pubkey: userPubkey})
	if err != nil {
		b.onError(fmt.Errorf("check user exists: %w", err))
		return "Hi. System error, failed to check your account"
	}
	if !isUserSaved {
		if err := b.dbConn.SaveUser(memory.User{Pubkey: userPubkey}); err != nil {
			b.onError(fmt.Errorf("create user: %w", err))
			return "Hi. System error, failed to create your account"
		}
	}

	return b.cfg.WelcomeMessage
}

func (b *uBot) onError(err error) {
	color.Red(err.Error())
}

func (b *uBot) onContactMessage(message structs.InstantMessage) {
	if len(message.Text) == channelIDLength { // TODO: check for hex
		channelID := message.Text

		// check channel exists
		isChannelSaved, err := b.dbConn.IsChannelExists(memory.Channel{
			ID: channelID,
		})
		if err != nil {
			// TODO
		}
		if isChannelSaved {
			// check ownership
			channelData, err := b.dbConn.GetChannel(channelID)
			if err != nil {
				// TODO
			}
			if channelData.OwnerPubkey != message.Pubkey {
				b.handler.SendContactMessage(
					message.Pubkey,
					"You have to be the owner of the channel",
				)
				return
			}
			b.handler.SendContactMessage(
				message.Pubkey,
				"TODO",
			)
			// TODO: go to command-enter mode
			return

		} else {
			if err := b.dbConn.SaveChannel(memory.Channel{
				ID:          channelID,
				OwnerPubkey: "", // TODO: get channel owner pubkey
			}); err != nil {
				// TODO
			}
		}

	}
}

func (b *uBot) onChannelMessage(message structs.WsChannelMessage) {
	// TODO
}

func (b *uBot) onPrivateChannelMessage(message structs.WsChannelMessage) {
	b.handler.SendChannelPrivateMessage(
		message.ChannelID,
		message.PubkeyHash,
		"Hi. Interested in this bot? "+
			"Check out the other sample projects:\n"+
			"https://udocs.gitbook.io/utopia-api/",
	)
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
