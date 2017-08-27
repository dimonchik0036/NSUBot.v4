package main

import (
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"strings"
)

var PagesMap Handlers

const (
	StrPageStart           = "start"
	StrPageError           = "error"
	StrPageErrorPermission = "error_permission"
	StrPage404NotFound     = "404_not_found"
	StrPageSuccess         = "success"
	StrPageWeather         = "weather" //2 button
	StrPageMenuMain        = "menu_main"
	StrPageMenuOption      = "menu_option"
	StrPageMenuSubscribers = "menu_subscribers"
	StrPageOptionLang      = "option_lang"
)

func initPagesMap() {
	PagesMap = NewHandlers(func(s string) string {
		s, _ = DecodePage(strings.ToLower(s))
		return s
	})
	PagesMap.AddHandler(Handler{Handler: PageStart}, StrPageStart)
	PagesMap.AddHandler(Handler{Handler: PageErrorPermission}, StrPageErrorPermission)
	PagesMap.AddHandler(Handler{Handler: Page404NotFound}, StrPage404NotFound)
	PagesMap.AddHandler(Handler{Handler: PageError}, StrPageError)
	PagesMap.AddHandler(Handler{Handler: PageWeather}, StrPageWeather)
	PagesMap.AddHandler(Handler{Handler: PageMenuMain}, StrPageMenuMain)
	PagesMap.AddHandler(Handler{Handler: PageMenuOption}, StrPageMenuOption)
	PagesMap.AddHandler(Handler{Handler: PageOptionLang}, StrPageOptionLang)
	PagesMap.AddHandler(Handler{Handler: PageMenuSubscribers}, StrPageMenuSubscribers)
	PagesMap.AddHandler(Handler{Handler: PageSuccess}, StrPageSuccess)
	initAdminPages(&PagesMap)
}

func DecodePage(p string) (string, string) {
	s := strings.SplitN(p, "*", 2)
	if len(s) > 1 {
		return s[0], s[1]
	}

	return s[0], ""
}

func PagesHandler(request *mapps.Request, subscriber *User) bool {
	var result string
	handler, ok := PagesMap.GetHandler(request.Page)
	if ok {
		if handler.PermissionLevel <= subscriber.Permission {
			result = handler.Handler(request, subscriber)
		} else {
			result = PageErrorPermission(request, subscriber)
		}
	} else {
		result = Page404NotFound(request, subscriber)
	}

	if result == "" {
		result = PageError(request, subscriber)
	}

	request.Ctx.WriteString(result)
	return true
}

func PageError(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageError, subscriber.Lang)
	return t.DoPage("")
}

func PageStart(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageStart, subscriber.Lang)
	return t.DoPage("")
}

func PageMenuMain(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageMenuMain, subscriber.Lang)
	return t.DoPage("")
}

func PageErrorPermission(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageErrorPermission, subscriber.Lang)
	return t.DoPage("")
}

func Page404NotFound(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPage404NotFound, subscriber.Lang)
	return t.DoPage("")
}

func PageSuccess(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageSuccess, subscriber.Lang)
	return t.DoPage("")
}

func PageWeather(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageWeather, subscriber.Lang)
	return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20)),
		mapps.Div("",
			mapps.Bold(GlobalWeather.ShowWeather())+mapps.Br+
				t.GetOptional(0)+GlobalWeather.ShowTime(),
		),
		t.Navigation,
	)
}

func PageMenuOption(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageMenuOption, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20)))
}

func PageOptionLang(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageOptionLang, subscriber.Lang)
	_, arg := DecodePage(request.Page)
	if arg == "" {
		return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20)))
	}

	subscriber.Lang = arg
	t = TextsForUsers.Get(StrPageOptionLang, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20)))
}

func PageMenuSubscribers(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageMenuSubscribers, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20)))
}
