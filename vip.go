package main

import (
	"TelegramBot/jokes"
	"github.com/dimonchik0036/Miniapps-pro-SDK"
)

const (
	CmdVipJoke = "joke"

	StrPageVipJoke = CmdVipJoke
)

func initVipCommands(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: CommandVipJoke, PermissionLevel: PermissionVIP}, CmdVipJoke)
}

func initVipPages(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: PageVipJoke, PermissionLevel: PermissionVIP}, StrPageVipJoke)
}

func CommandVipJoke(request *mapps.Request, subscriber *User) string {
	return PageVipJoke(request, subscriber)
}

func PageVipJoke(request *mapps.Request, subscriber *User) string {
	joke, err := jokes.GetJokes()
	if err != nil {
		joke = "Произошла ошибка, повторите попытку."
	}

	return mapps.Page("",
		mapps.Div("", mapps.EscapeString(joke)),
		mapps.Navigation("",
			mapps.Link("",
				StrPageVipJoke,
				"Новый анекдот",
			),
			mapps.Link("",
				StrPageMenuMain,
				"В меню",
			),
		))
}
