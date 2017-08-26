package main

import ()

func main() {
	println("Start")
	//initDefaultLog() //comment while testing
	initSystemLog()
	initConfig(LoadConfig())
	loadBotConfig()
	go UpdateSection(GlobalConfig, NewsHandler)
	HandlerStart()
}
