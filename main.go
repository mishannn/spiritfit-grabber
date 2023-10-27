package main

import (
	"log"
	"strconv"
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

	time_ := time.Now()
	fullness := int(clubDetails.Fullness * 100)

	err = SaveClubFullnessToSheet(cfg.GSheets.SheetID, cfg.GSheets.DataRange, time_, fullness, temp, feelsLike, windSpeed, rainLevel, snowLevel, pressure, humidity)
	if err != nil {
		log.Fatalf("Can't write club load to sheet: %s", err)
	}
}
