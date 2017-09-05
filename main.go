package main

func main() {
	println("Start")
	loadAllTexts(FilenameTextsUsers, &TextsForUsers)
	loadAllTexts(FilenameTextsAdmin, &TextsForAdmin)
	initDefaultLog() //comment while testing
	initSystemLog()
	initConfig(LoadConfig())
	loadBotConfig()
	initPagesMap()
	initCommandsMap()
	initVkSites()
	initBotNews()
	go UpdateSection(GlobalConfig, NewsHandler)
	HandlerStart()
}
