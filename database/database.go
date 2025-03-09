package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// database object
var BlockDB *gorm.DB

func InitDatabase() {
	// 1. We need GORM and db drivers
	// 2. Connect to the db using Gorm's Open Method
	var err error
	BlockDB, err = gorm.Open(sqlite.Open("blockchain.db"))
	if err != nil {
		panic("failed to connect to the database")
	}
	fmt.Println("Database connected!")
	// 3. Define models/structs for db schemas
	// 4. AutoMigrate - Create/update dB schemas based on our models
	err = BlockDB.AutoMigrate(&Block{}, &Transaction{}, &Withdrawal{}, &Account{})
	if err != nil {
		panic("failed to migrate")
	}
	fmt.Println("Database Migrated.")
}