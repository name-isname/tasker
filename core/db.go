package core

import (
	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
)

var DB *gorm.DB

// InitDB initializes the SQLite database
func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
