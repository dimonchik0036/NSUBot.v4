package main

import "github.com/dimonchik0036/Miniapps-pro-SDK"

const (
	CmdAdminReloadTexts = "reload_texts"
	CmdAdminMenu = "admin_menu"

	StrPageAdminReloadTexts = CmdAdminReloadTexts
)

func initAdminCommands(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: CommandGod}, CmdGod)
	handlers.AddHandler(Handler{Handler: CommandAdminReloadTexts, PermissionLevel: PermissionAdmin}, CmdAdminReloadTexts)
	handlers.AddHandler(Handler{Handler: CommandAdminMenu, PermissionLevel:PermissionAdmin}, CmdAdminMenu, "admin")
}

func initAdminPages(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler:PageAdminMenu, PermissionLevel:PermissionAdmin}, CmdAdminMenu)
	handlers.AddHandler(Handler{Handler:PageAdminReloadTexts,PermissionLevel:PermissionAdmin}, StrPageAdminReloadTexts)
}

func CommandAdminReloadTexts(request *mapps.Request, subscriber *User) string {
	_, args := DecodeCommand(request.Event.Text)
	request.Page += "*"+args
	return PageAdminReloadTexts(request, subscriber)
}

func PageAdminReloadTexts(request *mapps.Request, subscriber *User) string {
	t := TextsForAdmin.Get(StrPageAdminReloadTexts, subscriber.Lang)
	_, args := DecodePage(request.Page)
	if args == "" {
		return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20)))
	}

	switch args {
	case "admin":
		loadAllTexts(FilenameTextsAdmin, &TextsForAdmin)
	case "users":
		loadAllTexts(FilenameTextsUsers, &TextsForUsers)
	default:
		return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20)),
			mapps.Div("",
				t.GetOptional(1),
			),
			t.Navigation,
		)
	}

	return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20)),
				mapps.Div("",
				t.GetOptional(0),
				),
		t.Navigation,
	)
}

func CommandGod(request *mapps.Request, subscriber *User) string {
	if subscriber.User == Admin {
		subscriber.Permission = PermissionAdmin
		return PageSuccess(request, subscriber)
	}

	return PageErrorPermission(request, subscriber)
}

func CommandAdminMenu(request *mapps.Request, subscriber *User) string {
	return PageAdminMenu(request, subscriber)
}

func PageAdminMenu(request *mapps.Request, subscriber *User) string {
	t := TextsForAdmin.Get(CmdAdminMenu, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20)))
}