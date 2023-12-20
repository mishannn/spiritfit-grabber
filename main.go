package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func upMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("clickhouse"); err != nil {
		return fmt.Errorf("can't set dialect for migrations: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("can't up migrations: %w", err)
	}

	return nil
}

func runApplication() int {
	cfg, err := NewConfig("config.yaml")
	if err != nil {
		log.Printf("can't get config: %s", err)
		return 1
	}

	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{cfg.Database.Address},
		Auth: clickhouse.Auth{
			Database: cfg.Database.Database,
			Username: cfg.Database.Username,
			Password: cfg.Database.Password,
		},
	})
	err = upMigrations(db)
	if err != nil {
		log.Printf("can't up migrations: %s", err)
		return 1
	}

	clubDetails, err := GetClubDetails(cfg.Spirit.Token, cfg.Spirit.ClubID)
	if err != nil {
		log.Printf("can't get club details: %s", err)
		return 1
	}

	collectTime := time.Now()
	fullness := int(clubDetails.Fullness * 100)

	_, err = db.Exec("INSERT INTO club_fullness (DateTime, Fullness) VALUES ($1, $2)", collectTime, fullness)
	if err != nil {
		log.Printf("can't save club fullness: %s", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(runApplication())
}
