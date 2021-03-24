package db

import (
	"embed"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

var (
	//go:embed migrations/*
	migrations embed.FS
)

func up() error {
	srcDrv, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return err
	}

	db, err := DB.DB()
	if err != nil {
		return err
	}
	dbDrv, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("httpfs", srcDrv, "mysql", dbDrv)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
