package bot

import (
	"bot/pkg/filter"
	"bot/pkg/memory"
	"fmt"

	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/consts"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
)

const (
	channelIDLength    = 32
	defaultAccountName = "account.db"
)

const defaultPrivateMessage = "Hi. Interested in this bot? " +
	"Check out the other sample projects:\n" +
	"https://udocs.gitbook.io/utopia-api/"

type uBot struct {
	handler *uchatbot.ChatBot
	dbConn  memory.Memory
	cfg     UBotConfig
}

type UBotConfig struct {
	WelcomeMessage        string          `json:"welcomeMessage"`
	AccountName           string          `json:"accountName"`
	AutoChangeAccountName bool            `json:"autoChangeAccountName"`
	UtopiaConfig          utopiago.Config `json:"utopia"`
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

	if cfg.AutoChangeAccountName {
		return b, b.fixAccountName(cfg.AccountName)
	}

	return b, nil
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
	// check user exists
	isUserSaved, err := b.dbConn.IsUserExists(memory.User{Pubkey: message.Pubkey})
	if err != nil {
		b.onError(fmt.Errorf("check user exists: %w", err))
		return
	}
	if !isUserSaved {
		if err := b.dbConn.SaveUser(memory.User{Pubkey: message.Pubkey}); err != nil {
			b.onError(fmt.Errorf("save user: %w", err))
			return
		}
	}

	// get user data
	userPubkey := message.Pubkey
	u, err := b.dbConn.GetUser(userPubkey)
	if err != nil {
		b.onError(fmt.Errorf("get user data: %w", err))
		return
	}

	var maskError bool
	if u.EnterCommandMode {
		err = b.handleUserCommand(u, message.Text)
	} else {
		maskError, err = b.handleUserTextRequest(u, message.Text)
	}
	if err != nil {
		if !maskError {
			if _, messageErr := b.handler.GetClient().SendInstantMessage(
				userPubkey,
				err.Error(),
			); messageErr != nil {
				b.onError(fmt.Errorf("send to user errorChannelIDMustBeSent: %w", err))
				return
			}
			return
		}

		b.onError(fmt.Errorf("handle user request: %w", err))
		if _, messageErr := b.handler.GetClient().SendInstantMessage(
			userPubkey,
			errorNotifyDevelopers.Error(),
		); messageErr != nil {
			b.onError(fmt.Errorf("send handle response error to user: %w", err))
			return
		}
	}
}

func (b *uBot) handleUserCommand(u memory.User, msgText string) error {
	// TODO
	return nil
}

func (b *uBot) isJoinedToChannel(channelID string) (bool, error) {
	channels, err := b.handler.GetClient().GetChannels(structs.GetChannelsTask{
		SearchFilter: channelID,
		ChannelType:  consts.ChannelTypeJoined,
	})
	if err != nil {
		return false, fmt.Errorf("get channels joined: %w", err)
	}

	return len(channels) == 1, nil
}

func (b *uBot) handleUserTextRequest(
	u memory.User,
	channelID string,
) (maskError bool, handleErr error) {
	if !filter.NewChannelsFilter().Use(channelID) {
		return false, errorChannelIDMustBeSent
	}

	// check channel ownership
	channelData, err := b.handler.GetClient().GetChannelInfo(channelID)
	if err != nil {
		return true, fmt.Errorf("get channel data: %w", err)
	}

	if u.Pubkey != channelData.Owner {
		return false, fmt.Errorf(
			"you must be the owner of channel %q to control its filters",
			channelData.Title,
		)
	}

	// check channel saved
	isChannelSaved, err := b.dbConn.IsChannelExists(memory.Channel{ID: channelID})
	if err != nil {
		return true, fmt.Errorf("check channel exists: %w", err)
	}
	if !isChannelSaved {
		// save channel
		defFiltersJSON, err := getDefaultFiltersJSON()
		if err != nil {
			return true, fmt.Errorf("get default filters: %w", err)
		}

		if err := b.dbConn.SaveChannel(memory.Channel{
			ID:          channelID,
			OwnerPubkey: channelData.Owner,
			FiltersJSON: defFiltersJSON,
		}); err != nil {
			return true, fmt.Errorf("save channel: %w", err)
		}
	}

	// get channel config from db
	channelBotConfig, err := b.dbConn.GetChannel(channelID)
	if err != nil {
		return true, fmt.Errorf("get bot channel config: %w", err)
	}

	// check owner
	if channelBotConfig.OwnerPubkey != channelData.Owner {
		// channel owner was changed: save actual owner
		if err := b.dbConn.SetChannelOwner(channelID, channelData.Owner); err != nil {
			return true, fmt.Errorf("set channel owner: %w", err)
		}
	}

	isJoined, err := b.isJoinedToChannel(channelID)
	if err != nil {
		return true, fmt.Errorf("check channel joined: %w", err)
	}
	if !isJoined {
		if _, err := b.handler.GetClient().JoinChannel(channelID); err != nil {
			return true, fmt.Errorf("join to channel: %w", err)
		}
	}

	filters, err := channelBotConfig.GetFilters()
	if err != nil {
		return true, fmt.Errorf("parse channel filters: %w", err)
	}

	msg := "Send me the number of the selected option:\n\n" +
		getCommandsMessage(filters)
	if _, err := b.handler.GetClient().SendInstantMessage(u.Pubkey, msg); err != nil {
		return true, fmt.Errorf("send user commands: %w", err)
	}

	if err := b.dbConn.ToogleUserCommandMode(u.Pubkey, true); err != nil {
		return true, fmt.Errorf("toogle user command mode: %w", err)
	}
	return false, nil
}

func (b *uBot) onChannelMessage(message structs.WsChannelMessage) {
	// TODO: check channel is connected

	// TODO: check bot moderator rights

	// TODO: filter message

	// TODO: remove spam
}

func (b *uBot) onPrivateChannelMessage(message structs.WsChannelMessage) {
	b.handler.SendChannelPrivateMessage(
		message.ChannelID,
		message.PubkeyHash,
		defaultPrivateMessage,
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
