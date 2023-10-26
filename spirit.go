package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ClubDetailsResponse struct {
	Result      ClubDetails `json:"result"`
	UserMessage string      `json:"userMessage"`
	ErrorCode   int64       `json:"errorCode"`
}

type ClubDetails struct {
	ID                   string           `json:"id"`
	Title                string           `json:"title"`
	Subway               string           `json:"subway"`
	Latitude             string           `json:"latitude"`
	Longitude            string           `json:"longitude"`
	BackgroundImage      string           `json:"backgroundImage"`
	Fullness             float64          `json:"fullness"`
	Phone                string           `json:"phone"`
	Email                string           `json:"email"`
	Address              string           `json:"address"`
	HasSchedule          bool             `json:"hasSchedule"`
	SubscriptionCost     SubscriptionCost `json:"subscriptionCost"`
	SubscriptionDiscount int64            `json:"subscriptionDiscount"`
	CommonSquare         int64            `json:"commonSquare"`
	TrainingSquare       int64            `json:"trainingSquare"`
	ForemanEquipment     int64            `json:"foremanEquipment"`
	CardioExercicers     int64            `json:"cardioExercicers"`
	GroupLessons         int64            `json:"groupLessons"`
	Lockers              int64            `json:"lockers"`
	Showers              int64            `json:"showers"`
	WorkingTime          string           `json:"workingTime"`
	Gallery              []string         `json:"gallery"`
	Social               []Social         `json:"social"`
}

type Social struct {
	Type string `json:"type"`
	Link string `json:"link"`
}

type SubscriptionCost struct {
	Value    int64  `json:"value"`
	Currency string `json:"currency"`
}

func GetClubDetails(token string, club string) (*ClubDetails, error) {
	url := fmt.Sprintf("https://app.spiritfit.ru/Fitness/hs/mobile/clubs/%s", club)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	var data ClubDetailsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("can't parse response body: %w", err)
	}

	if data.ErrorCode != 0 {
		return nil, fmt.Errorf("server returned error: %s", data.UserMessage)
	}

	return &data.Result, nil
}
