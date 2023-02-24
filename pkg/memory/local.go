package memory

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type localDB struct {
	conn *gorm.DB
}

func NewLocalDB(filename string) (Memory, error) {
	fmt.Println("connect to db..")
	lg := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{
		Logger: lg,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	fmt.Println("migrate..")
	for _, prefab := range models {
		if err := db.AutoMigrate(prefab); err != nil {
			return nil, fmt.Errorf("failed to migrate: %w", err)
		}
	}

	return &localDB{
		conn: db,
	}, nil
}

func (db *localDB) isEntryExists(entryPointer interface{}, typePointer interface{}) (bool, error) {
	result := db.conn.Where(entryPointer).First(typePointer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (db *localDB) IsUserExists(u User) (bool, error) {
	return db.isEntryExists(&u, &User{})
}

func (db *localDB) SaveUser(u User) error {
	return db.conn.Save(&u).Error
}

func (db *localDB) IsChannelExists(c Channel) (bool, error) {
	return db.isEntryExists(&c, &Channel{})
}

func (db *localDB) SaveChannel(c Channel) error {
	return db.conn.Save(&c).Error
}

func (db *localDB) GetChannel(channelID string) (Channel, error) {
	c := Channel{}
	result := db.conn.Where(&Channel{
		ID: channelID,
	}).First(&c)
	return c, result.Error
}

func (db *localDB) ToogleUserCommandMode(pubkey string, enabled bool) error {
	return db.conn.Model(&User{}).Where("Pubkey", pubkey).
		Updates(User{
			EnterCommandMode: enabled,
		}).Error
}

func (db *localDB) SetUserPayload(u User, payload string) error {
	return db.conn.Model(&User{}).Where("Pubkey", u.Pubkey).
		Updates(User{
			Payload: payload,
		}).Error
}

func (db *localDB) GetUser(pubkey string) (User, error) {
	u := User{}
	result := db.conn.Where(&User{
		Pubkey: pubkey,
	}).First(&u)
	return u, result.Error
}

func (db *localDB) SetChannelOwner(channelID, ownerPubkey string) error {
	return db.conn.Model(&Channel{}).Where("OwnerPubkey", ownerPubkey).
		Updates(Channel{
			OwnerPubkey: ownerPubkey,
		}).Error
}
