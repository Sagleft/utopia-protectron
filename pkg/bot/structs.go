package bot

import (
	"bot/pkg/memory"

	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
)

type UBotConfig struct {
	WelcomeMessage        string          `json:"welcomeMessage"`
	AccountName           string          `json:"accountName"`
	AutoChangeAccountName bool            `json:"autoChangeAccountName"`
	UtopiaConfig          utopiago.Config `json:"utopia"`
}

// channel ID -> moderator rights
type channelModeratorRights map[string]structs.ModeratorRights

type channelFiltersData map[string]memory.UserFilters
