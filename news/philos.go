package news

import (
	"html"
	"regexp"
	"strings"
	"time"
)

const (
	PhilosHref       = "http://philos.nsu.ru/"
	PhilosTimeLayout = "02.01.2006"
	PhilosFuncName   = "philosname"
)

func PhilosNews() []*Site {
	return []*Site{
		&Site{
			Title:        "Новости",
			URL:          PhilosHref + "/left.htm",
			OptionalURL:  "/left.htm",
			NewsFunc:     Philos,
			NewsFuncName: PhilosFuncName,
		},
	}
}

func Philos(href string, count int) (news []News, err error) {
	body, err := getNewsPage(PhilosHref + href)
	if err != nil {
		return []News{}, err
	}

	rg, err := regexp.Compile("<ul class=\"myclass\">.*?</ul>")
	if err != nil {
		return []News{}, err
	}

	replaceRg, err := regexp.Compile("<!--.*?-->")
	if err != nil {
		return []News{}, err
	}

	body = replaceRg.ReplaceAll(body, []byte(""))
	replaceRg, err = regexp.Compile("</?font.*?>")
	if err != nil {
		return []News{}, err
	}

	body = replaceRg.ReplaceAll(rg.Find(body), []byte(""))

	rg, err = regexp.Compile("<li>.*?</li>")
	if err != nil {
		return []News{}, err
	}

	replaceRg, err = regexp.Compile("<.*?>")
	if err != nil {
		return []News{}, err
	}

	for _, b := range rg.FindAll(body, count) {
		href := hrefProcessing(b, 1)
		t, _ := time.ParseInLocation("<li>"+PhilosTimeLayout, string(b[:14]), time.Local)
		b = replaceRg.ReplaceAll(b, []byte(""))
		news = append(news, News{
			ID: func() int64 {
				h := string(href[0][0])
				if strings.Contains(h, "actual/") {
					return idScan(h)
				}
				return 0
			}(),
			Title: html.UnescapeString(string(b[11:])),
			URL:   PhilosHref + string(href[0][0]),
			Date:  t.Unix(),
		})
	}

	return
}
