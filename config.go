package main

import (
	"encoding/json"
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"github.com/dimonchik0036/NSUBot/news"
	"github.com/dimonchik0036/NSUBot/nsuschedule"
	"github.com/dimonchik0036/NSUBot/nsuweather"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

var Admin mapps.User
var SystemLoger *log.Logger
var GlobalWeather *nsuweather.Weather
var GlobalSites *Sites
var GlobalSchedule *nsuschedule.Schedule
var GlobalConfig *Config
var Port string

func initSystemLog() {
	file, err := os.OpenFile("syslog"+time.Now().Format("2006-01-02T15-04-05")+".txt", os.O_CREATE|os.O_WRONLY, os.FileMode(0700))
	if err != nil {
		log.Panic(err)
	}

	SystemLoger = log.New(file, "", log.LstdFlags)
}

func initDefaultLog() {
	file, err := os.OpenFile("log"+time.Now().Format("2006-01-02T15-04-05")+".txt", os.O_CREATE|os.O_WRONLY, os.FileMode(0700))
	if err != nil {
		log.Panic(err)
	}

	log.SetOutput(file)
}

func initConfig(config *Config) {
	config.Mux.Lock()
	defer config.Mux.Unlock()
	GlobalWeather = config.Weather
	GlobalSites = config.Sites
	GlobalSchedule = config.Schedule
	GlobalSubscribers = config.Users
	GlobalConfig = config
}

func loadBotConfig() {
	data, err := ioutil.ReadFile(".bot_config")
	if err != nil {
		log.Panicf("Bot config not found: %s", err.Error())
		return
	}

	tmp := struct {
		VkKey string
		Admin mapps.User
		Port  string
	}{}

	if err := yaml.Unmarshal(data, &tmp); err != nil {
		log.Panicf("Bot config: yaml throw error: %s", err.Error())
		return
	}

	news.SetVkServiceKey(tmp.VkKey)
	Admin = tmp.Admin
	Port = tmp.Port

	return
}

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
