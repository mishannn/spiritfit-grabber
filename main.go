package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ClubFullness struct {
	ID        int64     `gorm:"primarykey"`
	Timestamp time.Time `gorm:"index"`
	Fullness  int
}

func buildDSN(host string, port int, user, password, dbname string) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   dbname,
	}

	return u.String()
}

func runApplication() error {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.Parse()

	cfg, err := NewConfig(configPath)
	if err != nil {
		return fmt.Errorf("can't get config: %w", err)
	}

	// ---

	db, err := gorm.Open(
		postgres.Open(
			buildDSN(
				cfg.Database.Address,
				cfg.Database.Port,
				cfg.Database.Username,
				cfg.Database.Password,
				cfg.Database.Database,
			),
		),
		&gorm.Config{},
	)
	if err != nil {
		return fmt.Errorf("can't open database: %w", err)
	}

	err = db.AutoMigrate(&ClubFullness{})
	if err != nil {
		return fmt.Errorf("can't apply migrations: %w", err)
	}

	// ---

	c := cron.New(cron.WithSeconds())

	_, err = c.AddFunc(cfg.CronWithSeconds, func() {
		clubDetails, err := GetClubDetails(cfg.Spirit.Token, cfg.Spirit.ClubID)
		if err != nil {
			log.Printf("can't get club details: %sw\n", err)
			return
		}

		clubFullness := ClubFullness{
			Timestamp: time.Now(),
			Fullness:  int(clubDetails.Fullness * 100),
		}

		tx := db.Create(&clubFullness)
		if tx.Error != nil {
			log.Printf("can't save club fullness: %s\n", tx.Error)
			return
		}

		log.Printf("written club fullness: %+v\n", clubFullness)
	})
	if err != nil {
		return fmt.Errorf("can't add scheduled func: %w", err)
	}

	// ---

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	c.Start()
	log.Println("started")

	// ---

	<-ctx.Done()
	log.Println("stopping...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c.Stop()
	log.Println("stopped")

	<-shutdownCtx.Done()

	return nil
}

func main() {
	if err := runApplication(); err != nil {
		log.Fatalf("error: %s", err)
	}

	log.Println("exit")
}
