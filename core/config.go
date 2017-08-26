package core

import (
	"encoding/json"
	"github.com/dimonchik0036/NSUBot/nsuschedule"
	"github.com/dimonchik0036/NSUBot/nsuweather"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Config struct {
	Mux      sync.Mutex
	Schedule *nsuschedule.Schedule
	Weather  *nsuweather.Weather
	Sites    *Sites
	Users    *Users
}

func NewConfig() (config Config) {
	schedule := nsuschedule.NewSchedule()
	weather := nsuweather.NewWeather()
	sites := NewSites()
	config.Weather = &weather
	config.Schedule = &schedule
	config.Sites = &sites
	config.Users = &Users{}
	return
}

func LoadConfig() *Config {
	return &Config{
		Users:   loadUsers(),
		Weather: loadWeather(),
		Sites:   loadSites(),
	}
}

func (c *Config) Save() {
	c.Mux.Lock()
	defer c.Mux.Unlock()
	saveWeather(c.Weather)
	saveUsers(c.Users)
	saveSites(c.Sites)
}

func (c *Config) Reset() {
	saveWeather(c.Weather)
	saveUsers(c.Users)
	saveSites(c.Sites)
	c.Mux.Lock()
	log.Print("Выключаюсь")
	os.Exit(0)
}

func saveAndMarshal(filename string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0700))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func saveUsers(users *Users) {
	users.Mux.RLock()
	defer users.Mux.RUnlock()
	if err := saveAndMarshal("users.json", users); err != nil {
		log.Print(err)
	}
}

func saveWeather(weather *nsuweather.Weather) {
	weather.Mux.RLock()
	defer weather.Mux.RUnlock()
	if err := saveAndMarshal("weather.json", weather); err != nil {
		log.Print(err)
	}
}

func saveSites(sites *Sites) {
	if err := saveAndMarshal("sites.json", sites); err != nil {
		log.Print(err)
	}
}

func loadAndUnmarshal(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &v)
}

func loadUsers() *Users {
	var u Users
	if err := loadAndUnmarshal("users.json", &u); err != nil {
		log.Print(err)
		return &Users{}
	}

	return &u
}

func loadWeather() *nsuweather.Weather {
	var w nsuweather.Weather
	if err := loadAndUnmarshal("weather.json", &w); err != nil {
		log.Print(err)
		w = nsuweather.NewWeather()
		log.Print("Load default weather")
		return &nsuweather.Weather{}
	}

	return &w
}

func loadSites() *Sites {
	var s Sites
	if err := loadAndUnmarshal("sites.json", &s); err != nil {
		log.Print(err)
		s = NewSites()
		log.Print("Load default sites")
		return &s
	}

	for _, site := range s.Sites {
		site.Site.InitFunc()
	}

	return &s
}
