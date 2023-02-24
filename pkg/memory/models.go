package memory

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

var models = []interface{}{
	&User{}, &Channel{},
}

type User struct {
	gorm.Model

	Pubkey           string `gorm:"primaryKey"`
	EnterCommandMode bool   `gorm:"default:false;"`
	Payload          string
}

// filter tag -> enabled
type UserFilters map[string]bool

type Channel struct {
	gorm.Model

	ID          string
	OwnerPubkey string
	FiltersJSON string
}

func (User) TableName() string {
	return "users"
}

func (Channel) TableName() string {
	return "channels"
}

func (c Channel) GetFilters() (UserFilters, error) {
	var f UserFilters
	if err := json.Unmarshal([]byte(c.FiltersJSON), &f); err != nil {
		return f, fmt.Errorf("decode filters: %w", err)
	}
	return f, nil
}
