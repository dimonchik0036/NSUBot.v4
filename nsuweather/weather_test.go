package nsuweather

import (
	"regexp"
	"testing"
)

func TestGetWeather(t *testing.T) {
	w, err := GetWeather()
	if err != nil {
		t.Fatal(err)
	}

	check(t, w)
}

func TestNewWeather(t *testing.T) {
	w := NewWeather()
	check(t, w.Weather)
}

func check(t *testing.T, weather string) {
	reg, err := regexp.Compile("НГУ")
	if err != nil {
		t.Fatal(err)
	}

	if check := reg.FindString(weather); check == "" {
		t.Fatal("Не найден НГУ")
	}
}
