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
