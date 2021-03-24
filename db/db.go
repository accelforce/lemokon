package db

import (
	"fmt"
	"time"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/accelforce/lemokon/util"
)

var (
	DB *gorm.DB
)

func init() {
	cfg, err := driver.ParseDSN(util.RequireEnv("DB_URL"))
	if err != nil {
		panic(fmt.Errorf("DB_URL parse error: %w", err))
	}
	cfg.MultiStatements = true
	cfg.ParseTime = true
	cfg.Loc = time.Local
	dsn := cfg.FormatDSN()

	if err := connect(dsn); err != nil {
		err = nil

		name := cfg.DBName
		cfg.DBName = ""
		if err := connect(cfg.FormatDSN()); err != nil {
			panic(fmt.Errorf("attempted creating database and failed to connect to database: %w", err))
		}

		if err := DB.Exec(fmt.Sprintf("CREATE DATABASE %s;", name)).Error; err != nil {
			panic(fmt.Errorf("failed to create database: %w", err))
		}

		if err := connect(dsn); err != nil {
			panic(fmt.Errorf("successfully created database and failed to reconnect to database: %w", err))
		}
	}

	if err := up(); err != nil {
		panic(fmt.Errorf("error migrating database: %w", err))
	}
}

func connect(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}
	DB = db
	return nil
}
