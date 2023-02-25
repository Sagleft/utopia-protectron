package bot

import (
	"bot/pkg/filter"
	"bot/pkg/memory"
	"fmt"
	"strconv"

	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/consts"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
)

const (
	channelIDLength         = 32
	commandMessageMaxLength = 1
	defaultAccountName      = "account.db"
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

func (b *uBot) onContactMessage(msg structs.InstantMessage) {
	// check user exists
	isUserSaved, err := b.dbConn.IsUserExists(memory.User{Pubkey: msg.Pubkey})
	if err != nil {
		b.onError(fmt.Errorf("check user exists: %w", err))
		return
	}
	if !isUserSaved {
		if err := b.dbConn.SaveUser(memory.User{Pubkey: msg.Pubkey}); err != nil {
			b.onError(fmt.Errorf("save user: %w", err))
			return
		}
	}

	// get user data
	userPubkey := msg.Pubkey
	u, err := b.dbConn.GetUser(userPubkey)
	if err != nil {
		b.onError(fmt.Errorf("get user data: %w", err))
		return
	}

	var message string
	if u.EnterCommandMode {
		message, err = b.handleUserCommand(u, msg.Text)
	} else {
		message, err = b.handleUserTextRequest(u, msg.Text)
	}
	if err != nil {
		b.onError(fmt.Errorf("handle user request: %w", err))

		if _, messageErr := b.handler.GetClient().SendInstantMessage(
			userPubkey,
			errorNotifyDevelopers.Error(),
		); messageErr != nil {
			b.onError(fmt.Errorf("send handle response error to user: %w", err))
		}
		return
	}

	if message == "" {
		return
	}

	if _, err := b.handler.GetClient().SendInstantMessage(
		userPubkey,
		message,
	); err != nil {
		b.onError(fmt.Errorf("send message to user: %w", err))
		return
	}
}

func (b *uBot) handleUserCommand(u memory.User, msgText string) (string, error) {
	if len(msgText) > commandMessageMaxLength {
		return "You must enter the option number", nil
	}

	if msgText == "0" {
		if err := b.dbConn.ToggleUserCommandMode(u.Pubkey, false); err != nil {
			return "", fmt.Errorf("toggle user mode: %w", err)
		}
		if err := b.dbConn.SetUserPayload(u, ""); err != nil {
			b.onError(err)
		}

		return "OK", nil
	}

	// parse command code
	commandCode, err := strconv.ParseInt(msgText, 10, 32)
	if err != nil {
		return "", fmt.Errorf("parse command from user input %q: %w", msgText, err)
	}

	// verify command code
	filters := filter.GetFiltersArray()
	commandFilterIndex := int(commandCode - 1)

	if !(commandFilterIndex >= 0 && commandFilterIndex < len(filters)) {
		return "Incorrect command code, must be the number of one of the items", nil
	}

	filterTag := filters[commandFilterIndex].GetTag()

	// TODO: verify user payload
	channelID := u.Payload

	// get channel filters
	channelData, err := b.dbConn.GetChannel(channelID)
	if err != nil {
		return "", fmt.Errorf("get channel data: %w", err)
	}

	// get channel filters
	channelFilters, err := channelData.GetFilters()
	if err != nil {
		return "", fmt.Errorf("get channel filters: %w", err)
	}

	// toggle user filter
	enabled, isExists := channelFilters[filterTag]
	if !isExists {
		enabled = false
	}
	enabled = !enabled
	channelFilters[filterTag] = enabled

	if err := b.dbConn.UpdateChannelFilters(channelID, channelFilters); err != nil {
		return "", fmt.Errorf("update channel filters: %w", err)
	}

	return getCommandsMessage(channelFilters), nil
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

func (b *uBot) joinChannel(channelID string) error {
	isJoined, err := b.isJoinedToChannel(channelID)
	if err != nil {
		return fmt.Errorf("check channel joined: %w", err)
	}
	if !isJoined {
		if _, err := b.handler.GetClient().JoinChannel(channelID); err != nil {
			return fmt.Errorf("join to channel: %w", err)
		}
	}
	return nil
}

func (b *uBot) checkChannelOwner(
	channelID string,
	channelBotConfig memory.Channel,
	channelData structs.ChannelData,
) error {
	if channelBotConfig.OwnerPubkey == channelData.Owner {
		return nil
	}

	// channel owner was changed: save actual owner
	if err := b.dbConn.SetChannelOwner(channelID, channelData.Owner); err != nil {
		return fmt.Errorf("set channel owner: %w", err)
	}
	return nil
}

func (b *uBot) checkChannelSaved(
	channelID string,
	channelData structs.ChannelData,
) error {
	isChannelSaved, err := b.dbConn.IsChannelExists(memory.Channel{ID: channelID})
	if err != nil {
		return fmt.Errorf("check channel exists: %w", err)
	}
	if isChannelSaved {
		return nil
	}

	// save channel
	defFiltersJSON, err := getDefaultFiltersJSON()
	if err != nil {
		return fmt.Errorf("get default filters: %w", err)
	}

	if err := b.dbConn.SaveChannel(memory.Channel{
		ID:          channelID,
		OwnerPubkey: channelData.Owner,
		FiltersJSON: defFiltersJSON,
	}); err != nil {
		return fmt.Errorf("save channel: %w", err)
	}
	return nil
}

func (b *uBot) checkChannelOwnership(
	channelID string,
	u memory.User,
	channelData structs.ChannelData,
) (message string, handleErr error) {
	if u.Pubkey != channelData.Owner {
		return fmt.Sprintf(
			"you must be the owner of channel %q to control its filters",
			channelData.Title,
		), nil
	}
	return "", nil
}

func (b *uBot) handleUserTextRequest(
	u memory.User,
	channelID string,
) (responseMessage string, handleErr error) {
	if !filter.NewChannelsFilter().Use(channelID) {
		return "write me the channel ID, " +
			"anti-spam filters for which you need to configure", nil
	}

	channelData, err := b.handler.GetClient().GetChannelInfo(channelID)
	if err != nil {
		return "", fmt.Errorf("get channel data: %w", err)
	}

	message, err := b.checkChannelOwnership(channelID, u, channelData)
	if err != nil {
		return "", err
	}
	if message != "" {
		return message, nil
	}

	if err := b.checkChannelSaved(channelID, channelData); err != nil {
		return "", err
	}

	// get channel config from db
	channelBotConfig, err := b.dbConn.GetChannel(channelID)
	if err != nil {
		return "", fmt.Errorf("get bot channel config: %w", err)
	}

	if err := b.checkChannelOwner(channelID, channelBotConfig, channelData); err != nil {
		return "", err
	}

	if err := b.joinChannel(channelID); err != nil {
		return "", err
	}

	filters, err := channelBotConfig.GetFilters()
	if err != nil {
		return "", fmt.Errorf("parse channel filters: %w", err)
	}

	msg := "Send me the number of the selected option:\n\n" +
		getCommandsMessage(filters)

	if err := b.dbConn.SetUserPayload(u, channelID); err != nil {
		return "", fmt.Errorf("set user payload: %w", err)
	}

	if err := b.dbConn.ToggleUserCommandMode(u.Pubkey, true); err != nil {
		return "", fmt.Errorf("toogle user command mode: %w", err)
	}
	return msg, nil
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
