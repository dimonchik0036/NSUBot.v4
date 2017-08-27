package main

import (
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"strings"
)

var CommandsMap Handlers

const (
	CmdMenu   = "menu"
	CmdStart  = "start"
	CmdGod    = "god"
	CmdMyName = "myname"
)

func initCommandsMap() {
	CommandsMap = NewHandlers(func(s string) string {
		if strings.HasPrefix(s, "/") {
			s = s[1:]
		}

		s, _ = DecodeCommand(strings.ToLower(s))
		return s
	})

	CommandsMap.AddHandler(Handler{Handler: CommandMenu}, CmdMenu)
	CommandsMap.AddHandler(Handler{Handler: CommandStart}, CmdStart)
	CommandsMap.AddHandler(Handler{Handler: CommandMyName}, CmdMyName)
	initAdminCommands(&CommandsMap)
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
			result = PageErrorPermission(request, subscriber)
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

func CommandMenu(request *mapps.Request, subscriber *User) string {
	return PageMenuMain(request, subscriber)
}

func CommandStart(request *mapps.Request, subscriber *User) string {
	return PageStart(request, subscriber)
}

func CommandMyName(request *mapps.Request, subscriber *User) string {
	return PageMyName(request, subscriber)
}
