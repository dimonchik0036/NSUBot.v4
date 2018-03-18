package main

import (
	"github.com/dimonchik0036/Miniapps-wrapper"
	"github.com/dimonchik0036/NSUBot/nsuschedule"
	"strconv"
	"strings"
	"time"
)

const (
	CmdAdminReloadTexts      = "reload_texts"
	CmdAdminMenu             = "admin_menu"
	CmdAdminChangePermission = "change_perm"
	CmdAdminVkMenu           = "admin_vk"
	CmdAdminSiteSync         = "site_sync"
	CmdAdminBotReset         = "reset"
	CmdAdminBotSave          = "save"
	CmdAdminSendMessage      = "send"
	CmdAdminUsersCount       = "count_users"
	CmdAdminParity           = "parity"

	StrPageAdminReloadTexts      = CmdAdminReloadTexts
	StrPageAdminChangePermission = CmdAdminChangePermission
	StrPageAdminMenu             = CmdAdminMenu
	StrPageAdminVkMenu           = CmdAdminVkMenu
	StrPageAdminSiteSync         = CmdAdminSiteSync
	StrPageAdminSendMessage      = CmdAdminSendMessage
	StrPageAdminParity           = CmdAdminParity
)

func initAdminCommands(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: CommandGod}, CmdGod)
	handlers.AddHandler(Handler{Handler: CommandAdminMenu, PermissionLevel: PermissionAdmin}, CmdAdminMenu, "admin")
	handlers.AddHandler(Handler{Handler: CommandAdminParity, PermissionLevel: PermissionAdmin}, CmdAdminParity)
	handlers.AddHandler(Handler{Handler: CommandAdminVkMenu, PermissionLevel: PermissionAdmin}, CmdAdminVkMenu)
	handlers.AddHandler(Handler{Handler: CommandAdminBotSave, PermissionLevel: PermissionAdmin}, CmdAdminBotSave)
	handlers.AddHandler(Handler{Handler: CommandAdminBotReset, PermissionLevel: PermissionAdmin}, CmdAdminBotReset)
	handlers.AddHandler(Handler{Handler: CommandAdminSiteSync, PermissionLevel: PermissionAdmin}, CmdAdminSiteSync)
	handlers.AddHandler(Handler{Handler: CommandAdminUsersCount, PermissionLevel: PermissionAdmin}, CmdAdminUsersCount)
	handlers.AddHandler(Handler{Handler: CommandAdminReloadTexts, PermissionLevel: PermissionAdmin}, CmdAdminReloadTexts)
	handlers.AddHandler(Handler{Handler: CommandAdminSendMessage, PermissionLevel: PermissionAdmin}, CmdAdminSendMessage)
	handlers.AddHandler(Handler{Handler: CommandAdminChangePermission, PermissionLevel: PermissionAdmin}, CmdAdminChangePermission)

}

func initAdminPages(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: PageAdminMenu, PermissionLevel: PermissionAdmin}, StrPageAdminMenu)
	handlers.AddHandler(Handler{Handler: PageAdminVkMenu, PermissionLevel: PermissionAdmin}, StrPageAdminVkMenu)
	handlers.AddHandler(Handler{Handler: PageAdminSiteSync, PermissionLevel: PermissionAdmin}, StrPageAdminSiteSync)
	handlers.AddHandler(Handler{Handler: PageAdminReloadTexts, PermissionLevel: PermissionAdmin}, StrPageAdminReloadTexts)
	handlers.AddHandler(Handler{Handler: PageAdminSendMessage, PermissionLevel: PermissionAdmin}, StrPageAdminSendMessage)
	handlers.AddHandler(Handler{Handler: PageAdminChangePermission, PermissionLevel: PermissionAdmin}, StrPageAdminChangePermission)
}

func CommandAdminParity(request *mapps.Request, subscriber *User) string {
	nsuschedule.GlobalParity.Change()
	GlobalSchedule.SetParity(nsuschedule.GlobalParity.GetParity())
	return PageSuccess(request, subscriber)
}

func CommandAdminUsersCount(request *mapps.Request, subscriber *User) string {
	return mapps.Page("", mapps.Div("", mapps.Bold("Всего "+strconv.Itoa(GlobalSubscribers.Len())+" пользователей.")))
}

func CommandAdminSiteSync(request *mapps.Request, subscriber *User) string {
	return PageAdminSiteSync(request, subscriber)
}

func CommandAdminBotReset(request *mapps.Request, subscriber *User) string {
	go func() {
		time.Sleep(5 * time.Second)
		GlobalConfig.Reset()
	}()

	return PageSuccess(request, subscriber)
}

func CommandAdminBotSave(request *mapps.Request, subscriber *User) string {
	GlobalConfig.Save()
	return PageSuccess(request, subscriber)
}

func PageAdminSiteSync(request *mapps.Request, subscriber *User) string {
	for _, u := range GlobalSubscribers.Users {
		u.Sites = nil
	}

	GlobalSites.Mux.Lock()
	defer GlobalSites.Mux.Unlock()
	for _, s := range GlobalSites.Sites {
		if s.Subscribers == nil {
			continue
		}

		for _, u := range s.Subscribers.GetAll() {
			sub := GlobalSubscribers.User(u)
			if sub != nil {
				sub.Sub(s.Site.URL)
			}
		}
	}

	return PageSuccess(request, subscriber)
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

func CommandAdminSendMessage(request *mapps.Request, subscriber *User) string {
	_, args := DecodeCommand(request.Event.Text)
	if args != "" {
		request.SetField(CmdAdminSendMessage, args)
	}

	return PageAdminSendMessage(request, subscriber)
}

func PageAdminChangePermission(request *mapps.Request, subscriber *User) string {
	t := TextsForAdmin.Get(StrPageAdminChangePermission, subscriber.Lang)
	key := request.GetField(CmdAdminChangePermission)
	if key == "" {
		return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
	}

	SystemLogger.Print(subscriber.Key() + " использовал " + StrPageAdminChangePermission + "*" + key)

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

func PageAdminSendMessage(request *mapps.Request, subscriber *User) string {
	t := TextsForAdmin.Get(StrPageAdminSendMessage, subscriber.Lang)
	key := request.GetField(StrPageAdminSendMessage)
	if key == "" {
		return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
	}

	args := strings.SplitN(key, " ", 2)
	if len(args) < 2 {
		return pageInputHelp(t, 0)
	}

	u := GlobalSubscribers.User(args[0])
	if u == nil {
		return pageInputHelp(t, 1)
	}

	err := u.SendMessageBlock(mapps.Div("", mapps.Data(args[1])))
	if err != nil {
		return pageInputHelp(t, 2)
	}

	return pageInputHelp(t, 3)
}
