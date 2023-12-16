package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"strconv"
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

func main() {
	cfg, err := NewConfig("config.yaml")
	if err != nil {
		log.Fatalf("Can't get config: %s", err)
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
		log.Fatalf("can't up migrations: %s", err)
	}

	clubDetails, err := GetClubDetails(cfg.Spirit.Token, cfg.Spirit.ClubID)
	if err != nil {
		log.Fatalf("Can't get club details: %s", err)
	}

	lat, err := strconv.ParseFloat(clubDetails.Latitude, 64)
	if err != nil {
		log.Fatalf("Can't parse club latitude: %s", err)
	}

	lon, err := strconv.ParseFloat(clubDetails.Longitude, 64)
	if err != nil {
		log.Fatalf("Can't parse club longitude: %s", err)
	}

	weather, err := GetWeather(cfg.OpenWeather.APIKey, lat, lon)
	if err != nil {
		log.Fatalf("Can't get weather: %s", err)
	}

	temp := ConvertKelvinToCelsius(weather.Current.Temp)
	feelsLike := ConvertKelvinToCelsius(weather.Current.FeelsLike)
	windSpeed := weather.Current.WindSpeed
	rainLevel := weather.Current.Rain.The1H
	snowLevel := weather.Current.Snow.The1H
	pressure := weather.Current.Pressure
	humidity := weather.Current.Humidity

	collectTime := time.Now()
	fullness := int(clubDetails.Fullness * 100)

	_, err = db.Exec("INSERT INTO club_fullness (DateTime, Fullness, Temp, FeelsLike, WindSpeed, RainLevel, SnowLevel, Pressure, Humidity) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		collectTime, fullness, temp, feelsLike, windSpeed, rainLevel, snowLevel, pressure, humidity)
	if err != nil {
		log.Fatalf("Can't save club fullness: %s", err)
	}
}
