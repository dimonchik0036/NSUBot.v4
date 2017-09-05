package main

import (
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type BotNews_st struct {
	Date int64
	Text string
}

func (b *BotNews_st) String() string {
	return "Дата изменения: " + time.Unix(b.Date, 0).Format("15:04 02.01.2006") + mapps.Br + mapps.Br + b.Text
}

type ByDate []BotNews_st

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date > a[j].Date }

const (
	newsFilename          = "tgbotnews.yaml"
	CmdAdminAddBotNews    = "addbotnews"
	CmdAdminReloadBotNews = "reloadbotnews"
	CmdBotNewsList        = "bnlist"

	StrPageAdminAddBotNews    = CmdAdminAddBotNews
	StrPageAdminReloadBotNews = CmdAdminReloadBotNews
	StrPageAdminMenuBotNews   = "botnewsmenu"
	StrPageBotNewsList        = CmdBotNewsList
)

var BotNews struct {
	Mux  sync.RWMutex
	News []BotNews_st
}

func initBotNewsCommand(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: CommandAdminAddBotNews, PermissionLevel: PermissionAdmin}, CmdAdminAddBotNews)
	handlers.AddHandler(Handler{Handler: CommandAdminReloadBotNews, PermissionLevel: PermissionAdmin}, CmdAdminReloadBotNews)
	handlers.AddHandler(Handler{Handler: CommandBotNewsList}, CmdBotNewsList)
}

func initBotNewsPages(handlers *Handlers) {
	handlers.AddHandler(Handler{Handler: PageBotNewsList}, StrPageBotNewsList)
	handlers.AddHandler(Handler{Handler: PageAdminAddBotNews, PermissionLevel: PermissionAdmin}, StrPageAdminAddBotNews)
	handlers.AddHandler(Handler{Handler: PageAdminMenuBotNews, PermissionLevel: PermissionAdmin}, StrPageAdminMenuBotNews)
	handlers.AddHandler(Handler{Handler: PageAdminReloadBotNews, PermissionLevel: PermissionAdmin}, StrPageAdminReloadBotNews)

}

func initBotNews() {
	botNews, err := loadBotNews(newsFilename)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	BotNews.News = botNews
}

func loadBotNews(filename string) ([]BotNews_st, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Loading %s is failed. Err: %s", filename, err.Error())
		return []BotNews_st{}, err
	}

	var botNews []BotNews_st
	if err := yaml.Unmarshal(data, &botNews); err != nil {
		log.Printf("Loading %s is failed. Err: %s", filename, err.Error())
		return []BotNews_st{}, err
	}

	sort.Sort(ByDate(botNews))

	return botNews, nil
}

func saveBotNews(filename string, botNews []BotNews_st) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.FileMode(0700))
	if err != nil {
		log.Printf("Saving %s is failed. Err: %s", filename, err.Error())
		return err
	}

	data, err := yaml.Marshal(botNews)
	if err != nil {
		log.Printf("Saving %s is failed. Err: %s", filename, err.Error())
		return err
	}

	if _, err := file.Write(data); err != nil {
		log.Printf("Saving %s is failed. Err: %s", filename, err.Error())
		return err
	}

	return file.Close()
}

func appendBotNews(botNews []BotNews_st, newNews BotNews_st) []BotNews_st {
	botNews = append(botNews, newNews)
	sort.Sort(ByDate(botNews))
	if err := saveBotNews(newsFilename, botNews); err != nil {
		log.Printf("%s", err.Error())
	}

	return botNews
}

func CommandAdminReloadBotNews(request *mapps.Request, subscriber *User) string {
	return PageAdminReloadBotNews(request, subscriber)
}

func PageAdminReloadBotNews(request *mapps.Request, subscriber *User) string {
	botNews, err := loadBotNews(newsFilename)
	if err != nil {
		return PageError(request, subscriber)
	}

	BotNews.Mux.Lock()
	defer BotNews.Mux.Unlock()
	BotNews.News = botNews
	return PageSuccess(request, subscriber)
}

func PageAdminAddBotNews(request *mapps.Request, subscriber *User) string {
	t := TextsForAdmin.Get(StrPageAdminAddBotNews, subscriber.Lang)
	key := request.GetField(StrPageAdminAddBotNews)
	if key == "" {
		return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
	}

	BotNews.Mux.Lock()
	defer BotNews.Mux.Unlock()
	BotNews.News = appendBotNews(BotNews.News, BotNews_st{
		Date: time.Now().Unix(),
		Text: key,
	})

	return pageInputHelp(t, 0)
}

func CommandAdminAddBotNews(request *mapps.Request, subscriber *User) string {
	_, args := DecodeCommand(request.Event.Text)
	if args != "" {
		request.SetField(StrPageAdminAddBotNews, args)
	}

	return PageAdminAddBotNews(request, subscriber)
}

func PageAdminMenuBotNews(request *mapps.Request, subscriber *User) string {
	t := TextsForAdmin.Get(StrPageAdminMenuBotNews, subscriber.Lang)
	return t.DoPage(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(40)))
}

func CommandBotNewsList(request *mapps.Request, subscriber *User) string {
	return PageBotNewsList(request, subscriber)
}

func PageBotNewsList(request *mapps.Request, subscriber *User) string {
	BotNews.Mux.RLock()
	defer BotNews.Mux.RUnlock()

	t := TextsForUsers.Get(StrPageBotNewsList, subscriber.Lang)
	var pageNumber int
	_, arg := DecodePage(request.Page)
	pageNumber, _ = strconv.Atoi(arg)

	if pageNumber >= len(BotNews.News) {
		pageNumber = len(BotNews.News) - 1
	}

	if pageNumber < 0 {
		pageNumber = 0
	}

	var text string = t.GetOptional(0) + strconv.Itoa(pageNumber+1) + "/" + strconv.Itoa(len(BotNews.News)) + mapps.Br
	if len(BotNews.News) == 0 {
		text = t.GetOptional(1)
	} else {
		text += BotNews.News[pageNumber].String()
	}

	navigation := mapps.Navigation("",
		mapps.Link("",
			StrPageBotNewsList+"*"+strconv.Itoa(pageNumber-1),
			t.GetOptional(2),
		),
		mapps.Link("",
			StrPageBotNewsList+"*"+strconv.Itoa(pageNumber+1),
			t.GetOptional(3),
		))

	return mapps.Page(mapps.Attributes(mapps.TelegramLinksRealignmentThreshold(20), "telegram.links.realignment.enabled: true"),
		mapps.Div("",
			text,
		),
		navigation,
		t.Navigation,
	)
}
