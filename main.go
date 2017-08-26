package main

import (
	"github.com/dimonchik0036/NSUBot/core"
)

func main() {
	println("Start")
	//initDefaultLog() //comment while testing
	initSystemLog()
	initConfig(core.LoadConfig())
	loadBotConfig()
	go UpdateSection(GlobalConfig, NewsHandler)
	HandlerStart()
}
