package main

import (
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"github.com/dimonchik0036/NSUBot/news"
	"strconv"
	"strings"
)

var PagesMap Handlers

const (
	StrPageStart             = "start"
	StrPageError             = "error"
	StrPageErrorPermission   = "error_permission"
	StrPage404NotFound       = "404_not_found"
	StrPageSuccess           = "success"
	StrPageWeather           = "weather" //2 button
	StrPageMenuMain          = "menu_main"
	StrPageMenuOption        = "menu_option"
	StrPageMenuSubscribers   = "menu_subscribers"
	StrPageOptionLang        = "option_lang"
	StrPageMyName            = CmdMyName
	StrPageSubscriptionsList = "subscriptions_list"
	StrPageSiteList          = "site_list"
	StrPageFeedback          = "feedback"

	siteOnOnePage = 5
)

func initPagesMap() {
	PagesMap = NewHandlers(func(s string) string {
		s, _ = DecodePage(strings.ToLower(s))
		return s
	})

	PagesMap.AddHandler(Handler{Handler: PageStart}, StrPageStart)
	PagesMap.AddHandler(Handler{Handler: PageError}, StrPageError)
	PagesMap.AddHandler(Handler{Handler: PageWeather}, StrPageWeather)
	PagesMap.AddHandler(Handler{Handler: PageSuccess}, StrPageSuccess)
	PagesMap.AddHandler(Handler{Handler: PageFeedback}, StrPageFeedback)
	PagesMap.AddHandler(Handler{Handler: PageSiteList}, StrPageSiteList)
	PagesMap.AddHandler(Handler{Handler: PageMenuMain}, StrPageMenuMain)
	PagesMap.AddHandler(Handler{Handler: PageMenuOption}, StrPageMenuOption)
	PagesMap.AddHandler(Handler{Handler: PageOptionLang}, StrPageOptionLang)
	PagesMap.AddHandler(Handler{Handler: Page404NotFound}, StrPage404NotFound)
	PagesMap.AddHandler(Handler{Handler: PageErrorPermission}, StrPageErrorPermission)
	PagesMap.AddHandler(Handler{Handler: PageMenuSubscribers}, StrPageMenuSubscribers)
	PagesMap.AddHandler(Handler{Handler: PageSubscriptionsList}, StrPageSubscriptionsList)

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
		if handler.PermissionLevel > subscriber.Permission {
			result = PageErrorPermission(request, subscriber)
		} else {
			result = handler.Handler(request, subscriber)
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

func PageMyName(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageMyName, subscriber.Lang)
	return mapps.Page("",
		mapps.Div("",
			mapps.Bold(subscriber.String()),
		),
		t.Navigation,
	)
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
	return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)),
		mapps.Div("",
			mapps.Bold(GlobalWeather.ShowWeather())+mapps.Br+
				t.GetOptional(0)+GlobalWeather.ShowTime(),
		),
		t.Navigation,
	)
}

func PageFeedback(request *mapps.Request, subscriber *User) string {
	key := request.GetField(StrPageFeedback)
	if key != "" {
		Admin.SendMessage(subscriber.FullString(" ") + " оставил отзыв: " + key)
		return PageMenuOption(request, subscriber)
	}

	to := TextsForUsers.Get(StrPageMenuOption, subscriber.Lang)
	t := TextsForUsers.Get(StrPageFeedback, subscriber.Lang)
	return mapps.Page("",
		mapps.Div("",
			mapps.Input("submit", StrPageFeedback, t.Title),
		),
		mapps.Navigation(mapps.FormatAttr("id", "submit"),
			mapps.Link("", StrPageFeedback, "submit")),
		mapps.Navigation("",
			mapps.Link("",
				StrPageMenuOption, to.BackButton)),
	)
}

func PageMenuOption(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageMenuOption, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
}

func PageOptionLang(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageOptionLang, subscriber.Lang)
	_, arg := DecodePage(request.Page)
	if arg == "" {
		return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
	}

	subscriber.Lang = arg
	t = TextsForUsers.Get(StrPageOptionLang, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
}

func PageMenuSubscribers(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageMenuSubscribers, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
}

func PageSiteList(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageSiteList, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
}

func PageSubscriptionsList(request *mapps.Request, subscriber *User) string {
	t := TextsForUsers.Get(StrPageSubscriptionsList, subscriber.Lang)
	var siteNumber int
	var pageNumber int
	var args []string
	_, arg := DecodePage(request.Page)
	if arg != "" {
		args = strings.SplitN(arg, "_", 3)
		if len(args) > 1 {
			siteNumber, _ = strconv.Atoi(args[0])
			pageNumber, _ = strconv.Atoi(args[1])
		} else {
			siteNumber, _ = strconv.Atoi(args[0])
		}
	}

	var siteList []*news.Site
	if siteNumber == 5 {
		vkGroupSites.Mux.RLock()
		defer vkGroupSites.Mux.RUnlock()
		siteList = vkGroupSites.Groups
	} else {
		siteList = news.GetSite(siteNumber)
	}

	if pageNumber < 0 {
		pageNumber = 0
	}

	if pageNumber*siteOnOnePage > len(siteList) {
		pageNumber = len(siteList) / siteOnOnePage
	}

	if request.BadCommand != "" {
		args = strings.Split(request.BadCommand, " ")
		for _, arg := range args {
			subID, err := strconv.Atoi(arg)
			if err == nil {
				sub(subID-1, pageNumber, siteOnOnePage, siteList, subscriber)
			}
		}
	} else {
		if len(args) > 2 {
			subID, err := strconv.Atoi(args[2])
			if err == nil {
				sub(subID, pageNumber, siteOnOnePage, siteList, subscriber)
			}
		}
	}

	var navigation string
	for i, site := range siteList[pageNumber*siteOnOnePage:] {
		if i == siteOnOnePage {
			break
		}

		navigation +=
			mapps.Link("",
				"subscriptions_list*"+strconv.Itoa(siteNumber)+"_"+strconv.Itoa(pageNumber)+"_"+strconv.Itoa(i),
				checkSite(site.URL, subscriber)+site.Title,
			)

	}
	navigation = mapps.Navigation("", navigation)
	navigation += mapps.Navigation("",
		mapps.Link("",
			"subscriptions_list*"+strconv.Itoa(siteNumber)+"_"+strconv.Itoa(pageNumber-1),
			t.GetOptional(0),
		),
		mapps.Link("",
			"subscriptions_list*"+strconv.Itoa(siteNumber)+"_"+strconv.Itoa(pageNumber+1),
			t.GetOptional(1),
		),
		mapps.Link("",
			func() string {
				if siteNumber == 5 {
					return StrPageMenuSubscribers
				}
				return StrPageSiteList
			}(),
			t.GetOptional(2),
		),
	)
	return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)),
		mapps.Div("",
			t.Body+mapps.Br+
				mapps.Bold("",
					t.GetOptional(3)+strconv.Itoa(pageNumber+1)+"/"+strconv.Itoa(len(siteList)/siteOnOnePage+func() int {
						if len(siteList) == 0 {
							return 1
						} else if len(siteList)%siteOnOnePage == 0 {
							return 0
						} else {
							return 1
						}
					}())),
		),
		navigation,
	)
}

func checkSite(url string, subscriber *User) string {
	if subscriber.CheckSub(url) {
		return "☑️ "
	}

	return "❌"
}

func sub(id int, pageNumber int, size int, siteList []*news.Site, subscriber *User) {
	if id < 0 || id >= size {
		return
	}

	i := id + pageNumber*size
	if i < len(siteList) {
		SubscriptionManagement(siteList[i].URL, subscriber)
	}
}
func SubscriptionManagement(href string, subscriber *User) {
	GlobalSites.ChangeSub(href, subscriber)
	subscriber.ChangeSub(href)
}
