package main

import (
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"log"
	"strings"
)

var CommandsMap Handlers

const (
	CmdMenu     = "menu"
	CmdStart    = "start"
	CmdGod      = "god"
	CmdMyName   = "myname"
	CmdFeedback = StrPageFeedback
	CmdHelp     = StrPageHelp
)

func initCommandsMap() {
	CommandsMap = NewHandlers(func(s string) string {
		if strings.HasPrefix(s, "/") {
			s = s[1:]
		}

		s, _ = DecodeCommand(strings.ToLower(s))
		return s
	})

	CommandsMap.AddHandler(Handler{Handler: CommandHelp}, CmdHelp)
	CommandsMap.AddHandler(Handler{Handler: CommandMenu}, CmdMenu)
	CommandsMap.AddHandler(Handler{Handler: CommandStart}, CmdStart)
	CommandsMap.AddHandler(Handler{Handler: CommandMyName}, CmdMyName)
	CommandsMap.AddHandler(Handler{Handler: CommandFeedback}, CmdFeedback)
	CommandsMap.AddHandler(Handler{Handler: CommandScheduleToday}, "today", "t", "сегодня")
	CommandsMap.AddHandler(Handler{Handler: CommandScheduleTomorrow}, "tomorrow", "tm", "завтра")

	initAdminCommands(&CommandsMap)
	initBotNewsCommand(&CommandsMap)
	initVipCommands(&CommandsMap)
}

func DecodeCommand(cmd string) (string, string) {
	s := strings.SplitN(cmd, "*", 2)
	if len(s) > 1 {
		return s[0], s[1]
	}

	return s[0], ""
}

func CommandHandler(request *mapps.Request, subscriber *User) bool {
	var result string
	handler, ok := CommandsMap.GetHandler(request.Event.Text)
	if ok {
		if handler.PermissionLevel > subscriber.Permission {
			log.Print("Ошибка доступа у " + subscriber.String())
			return false
		} else {
			result = handler.Handler(request, subscriber)
		}
	}

	if result != "" {
		request.Ctx.WriteString(result)
		return true
	}

	return false
}

func CommandFeedback(request *mapps.Request, subscriber *User) string {
	_, args := DecodeCommand(request.Event.Text)
	if args != "" {
		request.SetField(StrPageFeedback, args)
	}

	return PageFeedback(request, subscriber)
}

func CommandMenu(request *mapps.Request, subscriber *User) string {
	return PageMenuMain(request, subscriber)
}

func CommandHelp(request *mapps.Request, subscriber *User) string {
	return PageHelp(request, subscriber)
}

func CommandStart(request *mapps.Request, subscriber *User) string {
	return PageStart(request, subscriber)
}

func CommandMyName(request *mapps.Request, subscriber *User) string {
	return PageMyName(request, subscriber)
}

func CommandScheduleToday(request *mapps.Request, subscriber *User) string {
	request.Page = StrPageShowSchedule + "*" + "0"
	return PageShowSchedule(request, subscriber)
}

func CommandScheduleTomorrow(request *mapps.Request, subscriber *User) string {
	request.Page = StrPageShowSchedule + "*" + "1"
	return PageShowSchedule(request, subscriber)
}
