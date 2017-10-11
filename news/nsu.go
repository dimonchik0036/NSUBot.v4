package news

import (
	"html"
	"regexp"
	"strconv"
	"time"
)

const (
	NsuHref       = "http://nsu.ru"
	NsuTimeLayout = "Mon, 2 Jan 2006 15:04:05 MST"
	NsuFuncName   = "nsu"
)

func NsuNews() []*Site {
	return []*Site{
		NsuMainPage(),
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
		URL:          NsuHref + "/fit?rss",
		OptionalURL:  "/fit?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuGGF() *Site {
	return &Site{
		Title:        "ГГФ",
		URL:          NsuHref + "/ggf?rss",
		OptionalURL:  "/ggf?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuIH() *Site {
	return &Site{
		Title:        "Гуманитарный институт",
		URL:          NsuHref + "/ih?rss",
		OptionalURL:  "/ih?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuHistory() *Site {
	return &Site{
		Title:        "История",
		URL:          NsuHref + "/ist?rss",
		OptionalURL:  "/ist?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuFundLing() *Site {
	return &Site{
		Title:        "Фундаментальная и прикладная лингвистика",
		URL:          NsuHref + "/fund_ling?rss",
		OptionalURL:  "/fund_ling?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuLing() *Site {
	return &Site{
		Title:        "Лингвистика",
		URL:          NsuHref + "/ling?rss",
		OptionalURL:  "/ling?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuJourn() *Site {
	return &Site{
		Title:        "Журналистика",
		URL:          NsuHref + "/journ?rss",
		OptionalURL:  "/journ?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuMed() *Site {
	return &Site{
		Title:        "ИМП (Здравоохранение)",
		URL:          NsuHref + "/med?rss",
		OptionalURL:  "/med?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuUF() *Site {
	return &Site{
		Title:        "ИФП (Юриспруденция)",
		URL:          NsuHref + "/uf?rss",
		OptionalURL:  "/uf?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}

}

func NsuPhilf() *Site {
	return &Site{
		Title:        "ИФП (Философия)",
		URL:          NsuHref + "/philf?rss",
		OptionalURL:  "/philf?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuMMF() *Site {
	return &Site{
		Title:        "ММФ",
		URL:          NsuHref + "/mmf?rss",
		OptionalURL:  "/mmf?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuFEN() *Site {
	return &Site{
		Title:        "ФЕН",
		URL:          NsuHref + "/fen?rss",
		OptionalURL:  "/fen?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuFF() *Site {
	return &Site{
		Title:        "ФФ",
		URL:          NsuHref + "/ff?rss",
		OptionalURL:  "/ff?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuEF() *Site {
	return &Site{
		Title:        "ЭФ",
		URL:          NsuHref + "/ef?rss",
		OptionalURL:  "/ef?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuMainPage() *Site {
	return &Site{
		Title:        "Все новости",
		URL:          NsuHref + "/news?rss",
		OptionalURL:  "/news?rss",
		NewsFunc:     NsuRss,
		NewsFuncName: NsuFuncName,
	}
}

func NsuRss(href string, count int) (news []News, err error) {
	body, err := getNewsPage(NsuHref + href)
	if err != nil {
		return []News{}, err
	}

	rg, err := regexp.Compile("<item>.*?</item>")
	if err != nil {
		return []News{}, err
	}

	hrefRg, err := regexp.Compile("<link>.*?</link>")
	if err != nil {
		return []News{}, err
	}

	guIdRg, err := regexp.Compile("<guid>.*?</guid>")
	if err != nil {
		return []News{}, err
	}

	titleRg, err := regexp.Compile("<title>.*?</title>")
	if err != nil {
		return []News{}, err
	}

	desRg, err := regexp.Compile("<description><![[]CDATA[[].*?]]></description>")
	if err != nil {
		return []News{}, err
	}

	repRg, err := regexp.Compile("<.*?>")
	if err != nil {
		return []News{}, err
	}

	help := func(start int, end int, body []byte) string {
		if len(body) > 0 {
			return string(body[start : len(body)-end])
		}
		return ""
	}

	for _, item := range rg.FindAll(body, count) {
		news = append(news, News{
			ID: func() int64 {
				id, _ := strconv.ParseInt(help(6, 7, guIdRg.Find(item)), 10, 64)
				return id
			}(),
			Title: html.UnescapeString(help(7, 8, titleRg.Find(item))),
			URL:   help(6, 7, hrefRg.Find(item)),
			Decryption: func() string {
				body := desRg.Find(item)
				body = body[22 : len(body)-17]
				body = repRg.ReplaceAll(body, []byte(""))
				var begin, end int
				for i, b := range body {
					if b != ' ' {
						if i > 0 {
							begin = i - 1
						}
						break
					}
				}

				for i := len(body) - 1; i >= 0; i-- {
					if body[i] != ' ' {
						if i < len(body) {
							end = i + 1
						}
						break
					}
				}

				return html.UnescapeString(string(body[begin:end]))
			}(),
			Date: 0,
		})
	}

	return
}

/*func Nsu(href string, count int) (news []News, err error) {
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

func NsuRss(href string, count int) (news []News, err error) {
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
}*/

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
