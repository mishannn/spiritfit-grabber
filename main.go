package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func sendToPrometheus(prometheusURL string, clubID string, count int) error {
	metricLine := fmt.Sprintf("spirit_fullness_percentage{club=\"%s\"} %d\n", clubID, count)

	resp, err := http.Post(prometheusURL, "text/plain", bytes.NewBufferString(metricLine))
	if err != nil {
		return fmt.Errorf("failed to send to prometheus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("prometheus returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func collectAndSendFullness(token, clubID, prometheusURL string) (int, error) {
	clubDetails, err := GetClubDetails(token, clubID)
	if err != nil {
		return 0, fmt.Errorf("can't get club details: %w", err)
	}

	clubFullness := int(clubDetails.Fullness * 100)
	sendToPrometheus(prometheusURL, clubID, clubFullness)

	return clubFullness, nil
}

func runApplication() error {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.Parse()

	cfg, err := NewConfig(configPath)
	if err != nil {
		return fmt.Errorf("can't get config: %w", err)
	}

	// first collect
	clubFullness, err := collectAndSendFullness(cfg.Spirit.Token, cfg.Spirit.ClubID, cfg.Metrics.PrometheusURL)
	if err != nil {
		log.Printf("can't collect and send fullness (1): %s\n", err)
	} else {
		log.Printf("written club fullness (1): %d%%\n", clubFullness)
	}

	// scheduled collect
	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc(cfg.CronWithSeconds, func() {
		clubFullness, err := collectAndSendFullness(cfg.Spirit.Token, cfg.Spirit.ClubID, cfg.Metrics.PrometheusURL)
		if err != nil {
			log.Printf("can't collect and send fullness (2): %s\n", err)
		} else {
			log.Printf("written club fullness (2): %d%%\n", clubFullness)
		}
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
