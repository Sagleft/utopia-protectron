package memory

import "gorm.io/gorm"

var models = []interface{}{
	&User{},
}

type User struct {
	gorm.Model

	Pubkey           string `gorm:"primaryKey"`
	EnterCommandMode bool   `gorm:"default:false;"`
	Payload          string
}

type Channel struct {
	gorm.Model

	ID          string
	OwnerPubkey string
	Filters     map[string]bool // filter tag -> enabled
}
