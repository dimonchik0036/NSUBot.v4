package news

import (
	"html"
	"regexp"
)

const (
	MmfHref       = "http://mmf.nsu.ru"
	MmfTimeLayout = "02.01.2006"
	MmfFuncName   = "mmfname"
)

func MmfNews() []*Site {
	return []*Site{
		&Site{
			Title:        "Новости",
			URL:          MmfHref + "/news/index",
			OptionalURL:  "/news/index",
			NewsFunc:     Mmf,
			NewsFuncName: MmfFuncName,
		},
		&Site{
			Title:        "Объявления",
			URL:          MmfHref + "/advert/index",
			OptionalURL:  "/advert/index",
			NewsFunc:     Mmf,
			NewsFuncName: MmfFuncName,
		},
		&Site{
			Title:        "Объявления студентам",
			URL:          MmfHref + "/students/advert",
			OptionalURL:  "/students/advert",
			NewsFunc:     Mmf,
			NewsFuncName: MmfFuncName,
		},
	}
}

func Mmf(href string, count int) (news []News, err error) {
	body, err := getNewsPage(MmfHref + href)
	if err != nil {
		return []News{}, err
	}

	rg, err := regexp.Compile("<div class=\"views-field views-field-title\">.*?</div>")
	if err != nil {
		return []News{}, err
	}

	dates := dateProcessing(body, count, "<span class=\"date-display-single\">", "</span>", MmfTimeLayout)

	for i, b := range rg.FindAll(body, count) {
		for _, v := range hrefProcessing(b, 1) {
			news = append(news, News{
				ID:    idScan(string(v[0])),
				Title: html.UnescapeString(string(v[1])),
				URL:   MmfHref + string(v[0]),
				Date:  dates[i],
			})
		}
	}

	return
}
