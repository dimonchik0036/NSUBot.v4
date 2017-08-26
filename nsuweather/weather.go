package nsuweather

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type Weather struct {
	Mux     sync.RWMutex `json:"-"`
	Weather string       `json:"weather"`
	Time    int64        `json:"time"`
}

const (
	WeatherLayout = "15:04 02.01.2006"
)

func GetWeather() (string, error) {
	res, err := http.Get("http://weather.nsu.ru/loadata.php")
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New("Status error: " + res.Status)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	reg, err := regexp.Compile("'Температура около .*?'")
	if err != nil {
		return "", err
	}

	currentWeather := reg.Find(body)
	if len(currentWeather) < 2 {
		return "", errors.New("Weather not found")
	}

	return string(currentWeather[1 : len(currentWeather)-1]), nil
}

func NewWeather() (weather Weather) {
	weather.Weather, _ = GetWeather()
	weather.Time = time.Now().Unix()
	return
}

func (weather *Weather) Update() {
	currentWeather, err := GetWeather()
	if err != nil {
		return
	}

	weather.Mux.Lock()
	defer weather.Mux.Unlock()
	weather.Weather = currentWeather
	weather.Time = time.Now().Unix()
}

func (weather *Weather) ShowWeather() (current string) {
	weather.Mux.RLock()
	defer weather.Mux.RUnlock()
	return weather.Weather
}

func (weather *Weather) ShowTime() string {
	weather.Mux.RLock()
	defer weather.Mux.RUnlock()
	return time.Unix(weather.Time, 0).Format(WeatherLayout)
}
