package gather

import (
	"encoding/json"
	"github.com/vibin18/openweather-exporter/signature"
	"log"
	"net/http"
	"strings"
)

const (
	baseUrl = "http://api.openweathermap.org/data/2.5/weather?"
)

type Sig struct {
	signature.Weather
}

type WeatherData interface {
	GetWeatherTemp() float64
	GetWeatherFeelsLike() float64
	GetWeatherPressure() float64
	GetWeatherHumidity() float64
	Fetchdata(cityName string, apiKey string)
}

func (m *Sig) GetWeatherTemp() float64 {
	return m.Main.Temp
}

func (m *Sig) GetWeatherFeelsLike() float64 {
	return m.Main.FeelsLike
}

func (m *Sig) GetWeatherPressure() float64 {
	return m.Main.Pressure
}

func (m *Sig) GetWeatherHumidity() float64 {
	return m.Main.Humidity
}

func (m *Sig) setWeatherHumidity(h float64) {
	m.Main.Humidity = h
}

func (m *Sig) setWeatherPressure(p float64) {
	m.Main.Pressure = p
}

func (m *Sig) setWeatherFeelsLike(f float64) {
	m.Main.FeelsLike = f
}

func (m *Sig) setWeatherTemp(t float64) {
	m.Main.Temp = t
}

func (m *Sig) Fetchdata(cityName string, apiKey string) {
	var weatherNow = *m
	key := strings.Trim(apiKey, "\"")
	log.Println("fetching weather data for",cityName)
	callUrlstr := []string{baseUrl, "q=", cityName, "&units=metric", "&appid=", key}
	callUrl := strings.Join(callUrlstr, "")
	res, err := http.Get(callUrl)
	if err != nil {
		log.Fatal(err)
	}

	json.NewDecoder(res.Body).Decode(&weatherNow)
	res.Body.Close()

	m.setWeatherFeelsLike(weatherNow.Main.FeelsLike)
	m.setWeatherHumidity(weatherNow.Main.Humidity)
	m.setWeatherPressure(weatherNow.Main.Pressure)
	m.setWeatherTemp(weatherNow.Main.Temp)

}
