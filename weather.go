package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Forecast struct {
	Query struct {
		Results struct {
			Channel struct {
				Item struct {
					Condition struct {
						Temp string `json:"temp"`
						Text string `json:"text"`
					} `json:"condition"`
				} `json:"item"`
			} `json:"channel"`
		} `json:"results"`
	} `json:"query"`
}

func GetForecast() (string, error) {
	resp, err := http.Get("https://query.yahooapis.com/v1/public/yql?q=select%20item.condition%20from%20weather.forecast%20where%20woeid%20in%20(select%20woeid%20from%20geo.places(1)%20where%20text%3D%22chicago%2C%20il%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys")
	if err != nil {
		return "", err
	}

	var forecast Forecast
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&forecast)

	condition := forecast.Query.Results.Channel.Item.Condition

	return fmt.Sprintf(
		"The weather today for Chicago is %s. The current temperature is %s degrees.",
		condition.Text,
		condition.Temp,
	), nil
}
