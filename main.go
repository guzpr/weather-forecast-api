package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CurrentWeather struct {
	Temperature float64 `json:"temperature"`
}
type WeatherRequest struct {
	CurrentWeather CurrentWeather `json:"current_weather"`
}

type WeatherResponse struct {
	TemperatureInCelcius    string `json:"temperature_in_celcius"`
	TemperatureInFahrenheit string `json:"temperature_in_fahrenheit"`
	Type                    string `json:"type"`
}

type Req struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
	router := gin.Default()
	router.GET("/weather", getWeather)

	router.Run("localhost:8080")
}
func getWeather(c *gin.Context) {
	var r Req
	if err := c.BindJSON(&r); err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	response, err := http.Get(fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true", r.Latitude, r.Longitude))

	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	var wreq WeatherRequest
	err = json.Unmarshal(responseData, &wreq)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	celcius := wreq.CurrentWeather.Temperature
	fahrenheit := (celcius * 1.8)

	celciusRounded := math.Round(celcius)
	var tipe string
	if celciusRounded <= 18 {
		tipe = "Cold"
	} else if celciusRounded >= 19 && celciusRounded <= 30 {
		tipe = "Warm"
	} else if celciusRounded >= 31 {
		tipe = "Hot"
	}

	c.JSON(http.StatusOK, WeatherResponse{
		TemperatureInCelcius:    fmt.Sprintf("%.2f°C", celcius),
		TemperatureInFahrenheit: fmt.Sprintf("%.2f°F", fahrenheit),
		Type:                    tipe,
	})
}
