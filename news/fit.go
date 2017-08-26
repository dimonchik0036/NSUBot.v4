package news

import (
	"html"
	"regexp"
	"strconv"
)

const (
	FitHref       = "http://fit.nsu.ru"
	FitTimeLayout = "02-01-06"
	FitFuncName   = "fitname"
)

func FitNews() []*Site {
	return []*Site{
		&Site{
			Title:        "События",
			URL:          FitHref + "/news/news",
			OptionalURL:  "/news/news",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
		&Site{
			Title:        "Объявления",
			URL:          FitHref + "/news/announc",
			OptionalURL:  "/news/announc",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
		&Site{
			Title:        "Конференции",
			URL:          FitHref + "/news/konf",
			OptionalURL:  "/news/konf",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
		&Site{
			Title:        "Конкурсы",
			URL:          FitHref + "/news/conc",
			OptionalURL:  "/news/conc",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
		&Site{
			Title:        "Вакансии",
			URL:          FitHref + "/news/vac",
			OptionalURL:  "/news/vac",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
		&Site{
			Title:        "Административные приказы",
			URL:          FitHref + "/news/administrativnye-prikazy",
			OptionalURL:  "/news/administrativnye-prikazy",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
	}
}

func FitChairs() []*Site {
	return []*Site{
		&Site{
			Title:        "Объявления кафедры систем информатики",
			URL:          FitHref + "/chairs/ksi/anksi",
			OptionalURL:  "/chairs/ksi/anksi",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
		&Site{
			Title:        "Объявления кафедры компьютерных систем",
			URL:          FitHref + "/chairs/kks/ankks",
			OptionalURL:  "/chairs/kks/ankks",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
		&Site{
			Title:        "Объявления кафедры общей информатики",
			URL:          FitHref + "/chairs/koi/koinews",
			OptionalURL:  "/chairs/koi/koinews",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
		&Site{
			Title:        "Объявления кафедры параллельных вычислений",
			URL:          FitHref + "/chairs/kpv/kpvnews",
			OptionalURL:  "/chairs/kpv/kpvnews",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
		&Site{
			Title:        "Объявления кафедры компьютерных технологий",
			URL:          FitHref + "/chairs/k-kt/kktnews",
			OptionalURL:  "/chairs/k-kt/kktnews",
			NewsFunc:     Fit,
			NewsFuncName: FitFuncName,
		},
	}
}

func Fit(href string, count int) (news []News, err error) {
	body, err := getNewsPage(FitHref + href + "?limit=" + strconv.Itoa(count))
	if err != nil {
		return []News{}, err
	}

	rg, err := regexp.Compile("<tbody>.*?</tbody>")
	if err != nil {
		return []News{}, err
	}

	body = rg.Find(body)
	dates := dateProcessing(body, count, "<td class=\"list-date\">", "</td>", FitTimeLayout)
	for i, v := range hrefProcessing(body, count) {
		news = append(news, News{
			ID:    idScan(string(v[0])),
			Title: html.UnescapeString(string(v[1])),
			URL:   FitHref + string(v[0]),
			Date:  dates[i],
		})
	}

	return
}

func FitNavigation(href string) ([]*Site, error) {
	body, err := getNewsPage(href)
	if err != nil {
		return []*Site{}, err
	}

	rg, err := regexp.Compile("<nav class=\"leftmenu\">.*?</nav>")
	if err != nil {
		return []*Site{}, err
	}

	return fitNavigationProcessing(rg.Find(body), Fit), nil
}

func fitNavigationProcessing(body []byte, f func(string, int) ([]News, error)) (sites []*Site) {
	for _, b := range hrefProcessing(body, -1) {
		sites = append(sites, &Site{
			OptionalURL: string(b[0]),
			Title:       string(b[1]),
			NewsFunc:    f,
		})
	}

	return
}
