package main

func main() {
	/*t := Texts{}
	t.Texts = map[string]*Localization{}
	t.Texts["menu_start"] = &Localization{
		Lang:map[string]string{
		LangRu:`Это огромный текст, который я сейчас скорпирую
		Так что пробуй!</br>`+"Приветствую!\n"+
			"Теперь я - ваш помощник.\n"+
			"Я позволяю получить быстрый доступ к температуре воздуха или же вы можете подписаться на рассылку новостей с различных сайтов и групп.\n"+
			"\n"+
			"Возможно будет полезным посмотреть /help, чтобы узнать все команды.\n"+
			"\n"+
			"При возникновении вопросов можно оставить /feedback или обратиться напрямую к @dimonchik0036.\n",
			LandEn: "Not russian button",
		},
	}
	t.Texts["help"] = &Localization{
		Lang:map[string]string{
			LangRu:"<div> Это<br/> хелп </div>",
			LandEn: "Not russian button",
		},
	}
	d, _ := yaml.Marshal(t)
	println(string(d))
	return*/
	println("Start")
	loadAllTexts(FilenameTextsUsers, &TextsForUsers)
	loadAllTexts(FilenameTextsAdmin, &TextsForAdmin)
	//initDefaultLog() //comment while testing
	initSystemLog()
	initConfig(LoadConfig())
	loadBotConfig()
	initPagesMap()
	initCommandsMap()
	initVkSites()
	go UpdateSection(GlobalConfig, NewsHandler)
	HandlerStart()
}
