package news

import (
	"html"
	"regexp"
)

const (
	FpHref       = "https://fp.nsu.ru"
	FpTimeLayout = "02.01.2006"
	FpFuncName   = "fpname"
)

func FpNews() []*Site {
	return []*Site{
		&Site{
			Title:        "Все новости",
			URL:          FpHref + "/content/news/",
			OptionalURL:  "/content/news/",
			NewsFunc:     Fp,
			NewsFuncName: FpFuncName,
		},
	}
}

func Fp(href string, count int) (news []News, err error) {
	body, err := getNewsPage(FpHref + href)
	if err != nil {
		return []News{}, err
	}

	rg, err := regexp.Compile("<div class=\"news-list\">.*?<div class=\"right_colomn\">")
	if err != nil {
		return []News{}, err
	}
	body = rg.Find(body)
	dates := dateProcessing(body, count, "<span class=\"news-date-time\">", "</span>", FpTimeLayout)
	hrefs := hrefProcessing(body, count*2)

	rg, err = regexp.Compile("</?b>")
	if err != nil {
		return []News{}, err
	}

	for i := range dates {
		news = append(news, News{
			ID:    idScan(string((hrefs[i*2+1][0])[16:])),
			Title: html.UnescapeString(string(rg.ReplaceAll(hrefs[i*2+1][1], []byte("")))),
			URL:   FpHref + string(hrefs[i*2+1][0]),
			Date:  dates[i],
		})
	}

	return
}
