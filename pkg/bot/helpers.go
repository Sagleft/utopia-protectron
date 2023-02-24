package bot

import (
	"bot/pkg/filter"
	"bot/pkg/memory"
	"encoding/json"
	"fmt"
)

func getEmojiTag(isFilterActivated bool) string {
	if isFilterActivated {
		return "[plus]"
	} else {
		return "[minus]"
	}
}

func getCommandsMessage(uFilters memory.UserFilters) string {
	msg := "[0] [bye-text]\n"

	filters := filter.GetFiltersArray()

	for i, f := range filters {

		isActivated, isFound := uFilters[f.GetTag()]
		if !isFound {
			isActivated = false
		}

		msg += fmt.Sprintf(
			"[%v] %s %s\n",
			i+1, getEmojiTag(isActivated), f.GetName(),
		)
	}

	return msg
}

func getDefaultFiltersJSON() (string, error) {
	f := memory.UserFilters{}

	for tag := range filter.GetFiltersMap() {
		f[tag] = false
	}

	fBytes, err := json.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("encode default filters: %w", err)
	}
	return string(fBytes), nil
}
