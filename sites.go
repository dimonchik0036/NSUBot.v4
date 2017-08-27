package main

import (
	"github.com/dimonchik0036/NSUBot/news"
	"log"
	"sync"
	"time"
)

type Site struct {
	Mux         sync.RWMutex `json:"-"`
	Site        *news.Site   `json:"site"`
	Subscribers *Set         `json:"users"`
}

func (s *Site) Sub(subscriber string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	if s.Subscribers == nil {
		s.Subscribers = new(Set)
	}

	s.Subscribers.Add(subscriber)
}

func (s *Site) UnSub(subscriber string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	if s.Subscribers == nil {
		s.Subscribers = new(Set)
	}
	s.Subscribers.Del(subscriber)
}

func (s *Site) ChangeSub(subscriber string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	if s.Subscribers == nil {
		s.Subscribers = new(Set)
	}

	s.Subscribers.Change(subscriber)
}

func (s *Site) Check(subscriber string) bool {
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	if s.Subscribers == nil {
		s.Subscribers = new(Set)
		return false
	}
	return s.Subscribers.Check(subscriber)
}

type Sites struct {
	Mux   sync.RWMutex     `json:"-"`
	Sites map[string]*Site `json:"sites"`
}

func NewSites() (sites Sites) {
	sites.Sites = map[string]*Site{}
	s := news.GetAllSites()
	for _, site := range s {
		sites.Sites[site.URL] = &Site{Site: site, Subscribers: &Set{Set: map[string]bool{}}}
	}
	return
}

func (s *Sites) AddSite(site *Site) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	if site == nil || site.Site == nil {
		log.Println("WTF?! Site is a nil pointer")
		return
	}

	s.Sites[site.Site.URL] = site
}

func (s *Sites) DelSite(href string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	delete(s.Sites, href)
}

func (s *Sites) ChangeSub(href string, user *User) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	site := s.Sites[href]
	if site == nil {
		log.Printf("%s wtf?!?! href %s not found", user.String(), href)
		return
	}
	site.ChangeSub(user.Key())
}

func (s *Sites) Sub(href string, user *User) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	site := s.Sites[href]
	if site == nil {
		log.Printf("%s wtf?!?! href %s not found", user.String(), href)
		return
	}

	site.Sub(user.Key())
}

func (s *Sites) UnSub(href string, user *User) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	site := s.Sites[href]
	if site == nil {
		log.Printf("%s wtf?!?! href %s not found", user.String(), href)
		return
	}

	site.UnSub(user.Key())
}

func (s *Sites) CheckUser(href string, user *User) bool {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	site := s.Sites[href]
	if site == nil {
		log.Printf("%s wtf?!?! href %s not found", user.String(), href)
		return false
	}

	return site.Check(user.Key())
}

func (s *Sites) Update(handler func([]string, []news.News, string)) {
	for _, site := range s.Sites {
		news, err := site.Site.Update(5)
		if err != nil {
			log.Printf("%s error: %s", site.Site.Title, err.Error())
			continue
		}

		if len(news) == 0 {
			continue
		}

		go handler(site.Subscribers.GetAll(), news, site.Site.Title)
		time.Sleep(250 * time.Millisecond)
	}
}
