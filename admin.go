package main

import (
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"strconv"
	"strings"
)

const (
	CmdAdminReloadTexts      = "reload_texts"
	CmdAdminMenu             = "admin_menu"
	CmdAdminChangePermission = "change_perm"
	CmdAdminVkMenu           = "admin_vk"

	StrPageAdminReloadTexts      = CmdAdminReloadTexts
	StrPageAdminChangePermission = CmdAdminChangePermission
	StrPageAdminMenu             = CmdAdminMenu
	StrPageAdminVkMenu           = CmdAdminVkMenu
)

func initAdminCommands(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: CommandGod}, CmdGod)
	handlers.AddHandler(Handler{Handler: CommandAdminMenu, PermissionLevel: PermissionAdmin}, CmdAdminMenu, "admin")
	handlers.AddHandler(Handler{Handler: CommandAdminVkMenu, PermissionLevel: PermissionAdmin}, CmdAdminVkMenu)
	handlers.AddHandler(Handler{Handler: CommandAdminReloadTexts, PermissionLevel: PermissionAdmin}, CmdAdminReloadTexts)
	handlers.AddHandler(Handler{Handler: CommandAdminChangePermission, PermissionLevel: PermissionAdmin}, CmdAdminChangePermission)
}

func initAdminPages(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: PageAdminMenu, PermissionLevel: PermissionAdmin}, StrPageAdminMenu)
	handlers.AddHandler(Handler{Handler: PageAdminVkMenu, PermissionLevel: PermissionAdmin}, StrPageAdminVkMenu)
	handlers.AddHandler(Handler{Handler: PageAdminReloadTexts, PermissionLevel: PermissionAdmin}, StrPageAdminReloadTexts)
	handlers.AddHandler(Handler{Handler: PageAdminChangePermission, PermissionLevel: PermissionAdmin}, StrPageAdminChangePermission)
}

func CommandAdminReloadTexts(request *mapps.Request, subscriber *User) string {
	_, args := DecodeCommand(request.Event.Text)
	request.Page += "*" + args
	return PageAdminReloadTexts(request, subscriber)
}

func PageAdminReloadTexts(request *mapps.Request, subscriber *User) string {
	t := TextsForAdmin.Get(StrPageAdminReloadTexts, subscriber.Lang)
	_, args := DecodePage(request.Page)
	if args == "" {
		return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
	}

	switch args {
	case "admin":
		loadAllTexts(FilenameTextsAdmin, &TextsForAdmin)
	case "users":
		loadAllTexts(FilenameTextsUsers, &TextsForUsers)
	default:
		return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)),
			mapps.Div("",
				t.GetOptional(1),
			),
			t.Navigation,
		)
	}

	return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)),
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
	t := TextsForAdmin.Get(StrPageAdminMenu, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
}

func CommandAdminChangePermission(request *mapps.Request, subscriber *User) string {
	_, args := DecodeCommand(request.Event.Text)
	if args != "" {
		request.SetField(CmdAdminChangePermission, args)
	}

	return PageAdminChangePermission(request, subscriber)
}

func PageAdminChangePermission(request *mapps.Request, subscriber *User) string {
	t := TextsForAdmin.Get(StrPageAdminChangePermission, subscriber.Lang)
	key := request.GetField(CmdAdminChangePermission)
	if key == "" {
		return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
	}

	SystemLoger.Print(subscriber.Key() + " использовал " + StrPageAdminChangePermission + "*" + key)

	if key == "reset all" && subscriber.Key() == Admin.Key() {
		for _, u := range GlobalSubscribers.Sample(func(u *User) bool {
			if u.Permission >= PermissionAdmin {
				return true
			}
			return false
		}) {
			u.Permission = PermissionVIP
		}
		return pageInputHelp(t, 4)
	}

	args := strings.SplitN(key, " ", 2)
	if len(args) < 2 {
		return pageInputHelp(t, 0)
	}

	perm, err := strconv.Atoi(args[0])
	if err != nil {
		return pageInputHelp(t, 1)
	}

	ok := GlobalSubscribers.ChangePermission(args[1], perm)
	if !ok {
		return pageInputHelp(t, 2)
	}

	return pageInputHelp(t, 3)
}

func pageInputHelp(t *Text, index int) string {
	return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)),
		mapps.Div("",
			mapps.Bold(
				t.GetOptional(index),
			),
		),
		t.Body,
		t.Navigation,
	)
}

func CommandAdminVkMenu(request *mapps.Request, subscriber *User) string {
	_, args := DecodeCommand(request.Event.Text)
	if args != "" {
		request.SetField(CmdAdminVkMenu, args)
	}

	return PageAdminVkMenu(request, subscriber)
}

func PageAdminVkMenu(request *mapps.Request, subscriber *User) string {
	t := TextsForAdmin.Get(StrPageAdminVkMenu, subscriber.Lang)
	key := request.GetField(StrPageAdminVkMenu)
	if key == "" {
		return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
	}

	args := strings.SplitN(key, " ", 2)
	if len(args) < 2 {
		return pageInputHelp(t, 0)
	}

	return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)),
		mapps.Div("",
			mapps.Bold(
				vkGroupProcessing(args[0], args[1]),
			),
		),
		t.Body,
		t.Navigation,
	)
}
