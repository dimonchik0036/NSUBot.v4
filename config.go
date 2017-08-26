package main

import (
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"github.com/dimonchik0036/NSUBot/core"
	"github.com/dimonchik0036/NSUBot/news"
	"github.com/dimonchik0036/NSUBot/nsuschedule"
	"github.com/dimonchik0036/NSUBot/nsuweather"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var Admin mapps.User
var SystemLoger *log.Logger
var GlobalWeather *nsuweather.Weather
var GlobalSites *core.Sites
var GlobalSchedule *nsuschedule.Schedule
var GlobalConfig *core.Config
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

func initConfig(config *core.Config) {
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
