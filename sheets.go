package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func SaveClubFullnessToSheet(spreadsheetId string, dataRange string, time time.Time, fullness int, temp float64, feelsLike float64, windSpeed float64, rainLevel float64, snowLevel float64) error {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("can't read credentials file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return fmt.Errorf("can't read JWT config from json: %w", err)
	}
	client := config.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("can't create sheets service: %w", err)
	}

	var vr sheets.ValueRange
	value := []interface{}{
		time.Format("01/02/2006 15:04:05"),
		fullness,
		temp,
		feelsLike,
		windSpeed,
		rainLevel,
		snowLevel,
	}
	vr.Values = append(vr.Values, value)

	_, err = srv.Spreadsheets.Values.Append(spreadsheetId, dataRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("can't write data to sheet: %w", err)
	}

	return nil
}
