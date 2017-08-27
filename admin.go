package main

import "github.com/dimonchik0036/Miniapps-pro-SDK"

const (
	CmdReloadTexts = "reload_texts"
)

func initAdminCommands(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: CommandGod}, CmdGod)
	handlers.AddHandler(Handler{Handler: CommandReloadTexts, PermissionLevel: PermissionAdmin}, CmdReloadTexts)
}

func CommandReloadTexts(request *mapps.Request, subscriber *User) string {
	_, args := DecodeCommand(request.Event.Text)
	switch args {
	case "admin":
		loadAllTexts(FilenameTextsAdmin, &TextsForUsers)
	case "users":
		loadAllTexts(FilenameTextsUsers, &TextsForUsers)
	default:
		return Page404NotFound(request, subscriber)
	}

	return PageSuccess(request, subscriber)
}

func CommandGod(request *mapps.Request, subscriber *User) string {
	if subscriber.User == Admin {
		subscriber.Permission = PermissionAdmin
		return PageSuccess(request, subscriber)
	}

	return PageErrorPermission(request, subscriber)
}
