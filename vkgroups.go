package main

import (
	"github.com/dimonchik0036/NSUBot/news"
	"sort"
	"strings"
	"sync"
)

type ByID []*news.Site

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

var vkGroupSites struct {
	Mux    sync.RWMutex
	Groups []*news.Site
}

func initVkSites() {
	GlobalSites.Mux.RLock()
	defer GlobalSites.Mux.RUnlock()

	var groups []*news.Site
	for key, site := range GlobalSites.Sites {
		if strings.HasPrefix(key, news.VkHref) {
			groups = append(groups, site.Site)
		}
	}

	sort.Sort(ByID(groups))
	vkGroupSites.Mux.Lock()
	defer vkGroupSites.Mux.Unlock()
	vkGroupSites.Groups = groups
}

func vkGroupProcessing(key string, data string) string {
	switch key {
	case "add":
		return addVkGroup(data)
	case "del":
		return delVkGroup(data)
	default:
		return "Неизвестная команда"
	}
}

func addVkGroup(data string) string {
	args := strings.SplitN(data, " ", 2)
	if len(args) < 2 {
		return "Недостаточно аргументов"
	}

	site := Site{
		Site: news.NewVkSite(int64(len(vkGroupSites.Groups)), args[0], args[1]),
	}

	_, err := site.Site.Update(2)
	if err != nil {
		return "Домен не найден"
	}

	vkGroupSites.Mux.Lock()
	defer vkGroupSites.Mux.Unlock()
	for _, s := range vkGroupSites.Groups {
		if s.OptionalURL == args[0] {
			return "Уже существует"
		}
	}
	vkGroupSites.Groups = append(vkGroupSites.Groups, site.Site)
	GlobalSites.AddSite(&site)
	return "Успешно"
}

func delVkGroup(data string) string {
	vkGroupSites.Mux.Lock()
	defer vkGroupSites.Mux.Unlock()
	var index int
	var groups []*news.Site
	var flag bool
	for _, site := range vkGroupSites.Groups {
		if site.OptionalURL == data {
			GlobalSites.DelSite(site.URL)
			flag = true
			break
		}
		groups = append(groups, site)
		index++
	}

	if !flag {
		return "Группа не найдена"
	}

	if index+1 < len(vkGroupSites.Groups) {
		for _, site := range vkGroupSites.Groups[index+1:] {
			site.ID--
			groups = append(groups, site)
		}
	}

	vkGroupSites.Groups = groups
	return "Успешно"
}
