package news

import (
	"errors"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type News struct {
	ID         int64  `json:"id"`
	URL        string `json:"url"`
	Title      string `json:"title"`
	Decryption string `json:"decryption"`
	Date       int64  `json:"date"`
}

const (
	TimeLayout = "02.01.2006"
)

func GetAllSites() (sites []*Site) {
	sites = append(sites, NsuNews()...)
	sites = append(sites, PhilosNews()...)
	sites = append(sites, FpNews()...)
	sites = append(sites, MmfNews()...)
	sites = append(sites, FitNews()...)
	sites = append(sites, FitChairs()...)

	return
}

const (
	FitNumber = iota
	NsuNumber
	MmfNumber
	FpNumber
	PhilosNumber
)

func GetSiteList() (list [][]*Site) {
	list = append(list, FitNews(), FitChairs())
	list = append(list, NsuNews())
	list = append(list, MmfNews())
	list = append(list, FpNews())
	list = append(list, PhilosNews())
	return list
}

func GetSite(number int) []*Site {
	switch number {
	case 0:
		return append(FitNews(), FitChairs()...)
	case 1:
		return NsuNews()
	case 2:
		return MmfNews()
	case 3:
		return FpNews()
	case 4:
		return PhilosNews()
	default:
		return []*Site{}
	}
}

type Site struct {
	ID            int64                                        `json:"id"`
	Mux           sync.Mutex                                   `json:"-"`
	Title         string                                       `json:"title"`
	OptionalTitle string                                       `json:"optional_title"`
	URL           string                                       `json:"url"`
	OptionalURL   string                                       `json:"option_url"`
	NewsFunc      func(href string, count int) ([]News, error) `json:"-"`
	NewsFuncName  string                                       `json:"news_func_name"`
	LastNews      News                                         `json:"last_news"`
}

func (s *Site) Update(countCheck int) (newNews []News, err error) {
	news, err := s.NewsFunc(s.OptionalURL, countCheck)
	if err != nil || len(news) == 0 {
		return newNews, err
	}

	s.Mux.Lock()
	defer s.Mux.Unlock()

	var max int64 = 0
	var maxIndex int = 0

	var maxID int64 = s.LastNews.ID
	var maxIDIndex int = -1
	for i, n := range news {
		if (s.LastNews.ID >= n.ID) && (n.ID != 0) || (n.URL == s.LastNews.URL) || s.LastNews.Date > n.Date {
			continue
		}

		if n.Date > max {
			max = n.Date
			maxIndex = i
		}

		if n.ID > maxID {
			maxID = n.ID
			maxIDIndex = i
		}

		if len(n.Decryption) > 3000 {
			n.Decryption = n.Decryption[:3000] + "...\nДля продолжения перейдите по ссылке в начале сообщения."
		}

		newNews = append(newNews, news[i])
	}

	if maxID != 0 {
		if maxID >= s.LastNews.ID && maxIDIndex != -1 {
			s.LastNews = news[maxIDIndex]
		}
	} else {
		if max >= s.LastNews.Date {
			s.LastNews = news[maxIndex]
		} else {
			s.LastNews = news[0]
		}
	}

	return reversNews(newNews), nil
}

func (s *Site) InitFunc() {
	switch s.NewsFuncName {
	case NsuFuncName:
		s.NewsFunc = NsuRss
	case FitFuncName:
		s.NewsFunc = Fit
	case PhilosFuncName:
		s.NewsFunc = Philos
	case MmfFuncName:
		s.NewsFunc = Mmf
	case FpFuncName:
		s.NewsFunc = Fp
	case VkFuncName:
		s.NewsFunc = Vk
	default:
		panic("WTF?!")
	}
}

func getNewsPage(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, errors.New("Error status: " + res.Status)
	}

	utf8, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		return []byte{}, err
	}

	rg, err := regexp.Compile("[\n\t]")
	if err != nil {
		return []byte{}, err
	}

	return rg.ReplaceAll(body, []byte("")), nil
}

func reversNews(news []News) (result []News) {
	for i := range news {
		result = append(result, news[len(news)-i-1])
	}

	return
}

func hrefProcessing(body []byte, count int) (result [][][]byte) {
	rg, err := regexp.Compile("<a.*?>.*?</a>")
	if err != nil {
		return
	}

	rgTitle, err := regexp.Compile(">")
	if err != nil {
		return
	}

	rgHref, err := regexp.Compile("\"")
	if err != nil {
		return
	}

	for _, href := range rg.FindAll(body, count) {
		href = href[9 : len(href)-4]
		titleInd := rgTitle.FindIndex(href)
		hrefInd := rgHref.FindIndex(href)
		result = append(result, [][]byte{href[:hrefInd[0]], href[titleInd[1]:]})
	}

	return
}

func idScan(url string) int64 {
	rg, err := regexp.Compile("[\\d]+")
	if err != nil {
		return 0
	}

	id, err := strconv.ParseInt(rg.FindString(url), 10, 64)
	if err != nil {
		return 0
	}

	return id
}

func dateProcessing(body []byte, count int, begin string, end string, layout string) (dates []int64) {
	rg, err := regexp.Compile(begin + ".*?" + end)
	if err != nil {
		return
	}

	for _, date := range rg.FindAll(body, count) {
		t, err := time.ParseInLocation(begin+layout+end, string(date), time.Local)
		if err != nil {
			dates = append(dates, 0)
		}
		dates = append(dates, t.Unix())
	}
	return
}
