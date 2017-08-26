package news

import (
	"html"
	"regexp"
	"time"
)

const (
	NsuHref        = "http://nsu.ru"
	NsuTimeLayout  = "02.01.2006 15:04"
	NsuFacFuncName = "nsufacname"
	NsuFuncName    = "nsuname"
)

func NsuNews() []*Site {
	return []*Site{
		NsuMainNews(),
		NsuMainPage(),
		NsuReportage(),
		NsuInterview(),
		NsuAnnounce(),
		NsuGGF(),
		NsuFIT(),
		NsuIH(),
		NsuHistory(),
		NsuFundLing(),
		NsuLing(),
		NsuJourn(),
		NsuMed(),
		NsuUF(),
		NsuPhilf(),
		NsuMMF(),
		NsuFEN(),
		NsuFF(),
		NsuEF(),
	}
}

func NsuFIT() *Site {
	return &Site{
		Title:        "ФИТ",
		URL:          NsuHref + "/fit",
		OptionalURL:  "/fit",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuGGF() *Site {
	return &Site{
		Title:        "ГГФ",
		URL:          NsuHref + "/ggf",
		OptionalURL:  "/ggf",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuIH() *Site {
	return &Site{
		Title:        "Гуманитарный институт",
		URL:          NsuHref + "/ih",
		OptionalURL:  "/ih",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuHistory() *Site {
	return &Site{
		Title:        "История",
		URL:          NsuHref + "/ist",
		OptionalURL:  "/ist",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuFundLing() *Site {
	return &Site{
		Title:        "Фундаментальная и прикладная лингвистика",
		URL:          NsuHref + "/fund_ling",
		OptionalURL:  "/fund_ling",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuLing() *Site {
	return &Site{
		Title:        "Лингвистика",
		URL:          NsuHref + "/ling",
		OptionalURL:  "/ling",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuJourn() *Site {
	return &Site{
		Title:        "Журналистика",
		URL:          NsuHref + "/journ",
		OptionalURL:  "/journ",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuMed() *Site {
	return &Site{
		Title:        "ИМП (Здравоохранение)",
		URL:          NsuHref + "/med",
		OptionalURL:  "/med",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuUF() *Site {
	return &Site{
		Title:        "ИФП (Юриспруденция)",
		URL:          NsuHref + "/uf",
		OptionalURL:  "/uf",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}

}

func NsuPhilf() *Site {
	return &Site{
		Title:        "ИФП (Философия)",
		URL:          NsuHref + "/philf",
		OptionalURL:  "/philf",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuMMF() *Site {
	return &Site{
		Title:        "ММФ",
		URL:          NsuHref + "/mmf",
		OptionalURL:  "/mmf",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuFEN() *Site {
	return &Site{
		Title:        "ФЕН",
		URL:          NsuHref + "/fen",
		OptionalURL:  "/fen",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuFF() *Site {
	return &Site{
		Title:        "ФФ",
		URL:          NsuHref + "/ff",
		OptionalURL:  "/ff",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuEF() *Site {
	return &Site{
		Title:        "ЭФ",
		URL:          NsuHref + "/ef",
		OptionalURL:  "/ef",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuAnnounce() *Site {
	return &Site{
		Title:        "События",
		URL:          NsuHref + "/news?mnc.news.type=announce",
		OptionalURL:  "/news?mnc.news.type=announce",
		NewsFunc:     Nsu,
		NewsFuncName: NsuFuncName,
	}
}

func NsuReportage() *Site {
	return &Site{
		Title:        "Репортажи",
		URL:          NsuHref + "/news?mnc.news.type=reportage",
		OptionalURL:  "/news?mnc.news.type=reportage",
		NewsFunc:     Nsu,
		NewsFuncName: NsuFuncName,
	}
}

func NsuMainPage() *Site {
	return &Site{
		Title:        "Все новости",
		URL:          NsuHref + "/news",
		OptionalURL:  "/news",
		NewsFunc:     Nsu,
		NewsFuncName: NsuFuncName,
	}
}

func NsuMainNews() *Site {
	return &Site{
		Title:        "Главные новости",
		URL:          NsuHref + "/?lang=ru",
		OptionalURL:  "/?lang=ru",
		NewsFunc:     NsuFac,
		NewsFuncName: NsuFacFuncName,
	}
}

func NsuInterview() *Site {
	return &Site{
		Title:        "Интервью",
		URL:          NsuHref + "/news?mnc.news.type=interview",
		OptionalURL:  "/news?mnc.news.type=interview",
		NewsFunc:     Nsu,
		NewsFuncName: NsuFuncName,
	}
}

func Nsu(href string, count int) (news []News, err error) {
	body, err := getNewsPage(NsuHref + href)
	if err != nil {
		return []News{}, err
	}

	rg, err := regexp.Compile("<div id=\"news-container\" class=\"list-holder\">.*?<div class=\"partners-holder\">")
	if err != nil {
		return []News{}, err
	}

	body = rg.Find(body)
	rg, err = regexp.Compile("<h.>.*?</p>")
	if err != nil {
		return []News{}, err
	}

	decryptionRg, err := regexp.Compile("<p>.*?</p>")
	if err != nil {
		return []News{}, err
	}

	for _, b := range rg.FindAll(body, count) {
		href := hrefProcessing(b, 1)
		news = append(news, News{
			Title: html.UnescapeString(string(href[0][1])),
			URL:   NsuHref + string(href[0][0]),
			Decryption: html.UnescapeString(func() string {
				if s := decryptionRg.Find(b); len(s) > 7 {
					return string(s[3 : len(s)-4])
				}
				return ""
			}()),
			Date: nsuDate(NsuHref + string(href[0][0])).Unix(),
		})

	}

	return
}

func NsuFac(href string, count int) (news []News, err error) {
	body, err := getNewsPage(NsuHref + href)
	if err != nil {
		return []News{}, err
	}

	rg, err := regexp.Compile("<ul class=\"news-list.*?</ul>")
	if err != nil {
		return []News{}, err
	}

	body = rg.Find(body)
	rg, err = regexp.Compile("<li.*?</li>")
	if err != nil {
		return []News{}, err
	}

	titleRg, err := regexp.Compile("<div class=\"text-holder\">.*?</div>")
	if err != nil {
		return []News{}, err
	}

	decryptionRg, err := regexp.Compile("<p>.*?</p>")
	if err != nil {
		return []News{}, err
	}

	for _, b := range rg.FindAll(body, count) {
		b = titleRg.Find(b)
		href := hrefProcessing(b, 1)
		news = append(news, News{
			Title: html.UnescapeString(string(href[0][1])),
			URL:   NsuHref + string(href[0][0]),
			Decryption: html.UnescapeString(func() string {
				if s := decryptionRg.Find(b); len(s) > 7 {
					return string(s[3 : len(s)-4])
				}
				return ""
			}()),
			Date: nsuDate(NsuHref + string(href[0][0])).Unix(),
		})

	}

	return
}

func nsuDate(url string) time.Time {
	body, err := getNewsPage(url)
	if err != nil {
		return time.Time{}
	}
	begin := "<p style=\"clear: both; text-align: right; color: grey; font-style: italic; font-size: small; margin-top: 10px;\">"
	end := "</p>"

	timeRg, err := regexp.Compile(begin + ".*?" + end)
	if err != nil {
		return time.Time{}
	}
	str := timeRg.Find(body)
	if len(str) < 165 {
		return time.Time{}
	}

	t, _ := time.ParseInLocation(NsuTimeLayout, string(str[165:len(str)-16]), time.Local)
	return t
}
