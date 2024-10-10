package database

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/carp-cobain/referrals/database/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectAndMigrate connects to a database and runs migrations using project models.
func ConnectAndMigrate() (*gorm.DB, *gorm.DB, error) {
	dsn := dsnEnvLookup()
	writeDB, err := Connect(dsn, 1)
	if err != nil {
		return nil, nil, err
	}
	readDB, err := Connect(dsn, max(4, runtime.NumCPU()))
	if err != nil {
		return nil, nil, err
	}
	if err := RunMigrations(writeDB); err != nil {
		return nil, nil, err
	}
	return readDB, writeDB, nil
}

// Connect to a sqlite3 database.
func Connect(dsn string, maxConns int) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Discard, //Default.LogMode(logger.Info),
	}
	db, err := gorm.Open(sqlite.Open(dsn), config)
	if err != nil {
		return nil, err
	}
	if err = setPragmas(db); err != nil {
		log.Printf("unable to set PRAGMAs for sqlite connection: %+v", err)
	}
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(maxConns)
	}
	return db, nil
}

// Run migrations on a database using project models.
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&model.Campaign{}, &model.Signup{})
}

// Optimize a sqlite database for production.
func setPragmas(db *gorm.DB) error {
	stmts := []string{
		"journal_mode = WAL",
		"busy_timeout = 5000",
		"synchronous = NORMAL",
		"cache_size = 1000000000",
		"foreign_keys = true",
		"temp_store = memory",
		"wal_autocheckpoint = 1000",
		"wal_checkpoint(TRUNCATE)",
	}
	for _, stmt := range stmts {
		if err := db.Exec(fmt.Sprintf("PRAGMA %s;", stmt)).Error; err != nil {
			return err
		}
	}
	return nil
}

// Lookup db dsn param from env var
func dsnEnvLookup() string {
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		log.Panicf("DB_DSN not defined")
	}
	return dsn
}
