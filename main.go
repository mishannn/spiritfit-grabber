package main

import (
	"log"
	"time"
)

func main() {
	cfg, err := NewConfig("config.yaml")
	if err != nil {
		log.Fatalf("Can't get config: %s", err)
	}

	clubDetails, err := GetClubDetails(cfg.Spirit.Token, cfg.Spirit.ClubID)
	if err != nil {
		log.Fatalf("Can't get club details: %s", err)
	}

	t := time.Now()
	l := int(clubDetails.Fullness * 100)

	err = SaveClubFullnessToSheet(cfg.GSheets.SheetID, cfg.GSheets.DataRange, t, l)
	if err != nil {
		log.Fatalf("Can't write club load to sheet: %s", err)
	}
}
