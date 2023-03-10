package memory

type Memory interface {
	// IsUserExists - check is user exists
	IsUserExists(u User) (bool, error)

	// GetUser data
	GetUser(pubkey string) (User, error)

	// SaveUser - save user entry
	SaveUser(u User) error

	// ToggleUserCommandMode - mark whether the user is expected to enter a command
	ToggleUserCommandMode(pubkey string, enabled bool) error

	// SetUserPayload - update the data interacting with via command mode
	SetUserPayload(u User, payload string) error

	IsChannelExists(c Channel) (bool, error)
	SaveChannel(c Channel) error
	GetChannel(channelID string) (Channel, error)
	GetChannels() ([]Channel, error)
	SetChannelOwner(channelID, ownerPubkey string) error
	UpdateChannelFilters(channelID string, f UserFilters) error
}
