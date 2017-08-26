package main

import (
	"github.com/dimonchik0036/NSUBot/core"
	"github.com/dimonchik0036/NSUBot/news"
	"github.com/dimonchik0036/NSUBot/nsuweather"
	"time"
)

func UpdateSection(config *core.Config, newsHandler func([]string, []news.News, string)) {
	go weatherUpdate(config.Weather, 2*time.Minute)

	go save(config, 20*time.Second, 5*time.Minute)

	go sitesUpdate(config.Sites, 45*time.Second, newsHandler)
}

func weatherUpdate(weather *nsuweather.Weather, duration time.Duration) {
	for {
		weather.Update()
		time.Sleep(duration)
	}
}

func save(config *core.Config, delay time.Duration, duration time.Duration) {
	time.Sleep(delay)
	for {
		config.Save()
		time.Sleep(duration)
	}
}

func sitesUpdate(sites *core.Sites, duration time.Duration, handler func([]string, []news.News, string)) {
	for {
		sites.Update(handler)
		time.Sleep(duration)
	}
}
