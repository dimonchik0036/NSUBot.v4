package news

import (
	"github.com/dimonchik0036/vk-api"
	"sort"
	"strings"
)

const (
	VkHref       = "http://vk.com"
	VkTimeLayout = "02-01-06"
	VkFuncName   = "vkname"
)

var vkClient *vkapi.Client

type ByID []News

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID > a[j].ID }

func Vk(domain string, count int) (news []News, err error) {
	_, wall, _, _, e := vkClient.GetWall(vkapi.NewDstFromDomain(domain), int64(count), 0, "all", false)
	if e != nil {
		return []News{}, e.ToError()
	}

	for _, w := range wall {
		news = append(news, News{
			ID:         w.ID,
			URL:        w.URL(),
			Date:       w.Date,
			Decryption: w.Text,
		})
	}

	sort.Sort(ByID(news))

	return
}

func SetVkServiceKey(key string) {
	vkClient, _ = vkapi.NewClientFromToken(key)
	vkClient.SetLanguage(vkapi.LangRU)
}

func NewVkSite(id int64, domain string, title string) *Site {
	return &Site{
		ID:           id,
		URL:          VkHref + "/" + strings.ToLower(domain),
		Title:        title,
		OptionalURL:  domain,
		NewsFunc:     Vk,
		NewsFuncName: VkFuncName,
	}
}
