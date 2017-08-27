package main

import (
	"fmt"
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"strconv"
	"sync"
	"time"
)

var GlobalSubscribers *Users

func CheckNewSubscriber(request *mapps.Request) *User {
	key := request.User()
	user := GlobalSubscribers.User(key.Key())
	if user == nil {
		return newUser(request)
	}
	return user
}

func newUser(request *mapps.Request) (subscriber *User) {
	subscriber = new(User)
	subscriber.Queue.Lock()
	defer subscriber.Queue.Unlock()
	subscriber.User = request.User()
	GlobalSubscribers.Add(subscriber)

	subscriber.Permission = PermissionUser
	subscriber.Lang = request.Lang
	t := time.Now().Unix()
	subscriber.DateCreated = t
	subscriber.DateLastActivities = t

	go Admin.SendMessageBlock(mapps.Div("", subscriber.NewUserString()))
	return
}

func refreshSubscriber(subscriber *User) *User {
	subscriber.DateLastActivities = time.Now().Unix()
	return subscriber
}

const (
	UserLayout = "2006/01/02 15:04:05"
)

type User struct {
	mapps.User
	SitesMux           sync.RWMutex `json:"-"`
	Sites              *Set         `json:"sites,omitempty"`
	Permission         int          `json:"permission"`
	DateCreated        int64        `json:"date_created"`
	DateLastActivities int64        `json:"date_last_activities"`
	MessageCount       int64        `json:"message_count"`
	Lang               string       `json:"lang"`
	Queue              sync.Mutex   `json:"-"`
}

func (u *User) Sub(href string) {
	u.SitesMux.Lock()
	defer u.SitesMux.Unlock()
	if u.Sites == nil {
		u.Sites = new(Set)
	}

	u.Sites.Add(href)
}

func (u *User) CheckSub(href string) bool {
	u.SitesMux.Lock()
	defer u.SitesMux.Unlock()
	if u.Sites == nil {
		u.Sites = new(Set)
		return false
	}

	return u.Sites.Check(href)
}

func (u *User) UnSub(href string) {
	u.SitesMux.Lock()
	defer u.SitesMux.Unlock()
	if u.Sites == nil {
		u.Sites = new(Set)
	}

	u.Sites.Del(href)
}

func (u *User) ChangeSub(href string) {
	u.SitesMux.Lock()
	defer u.SitesMux.Unlock()
	if u.Sites == nil {
		u.Sites = new(Set)
	}

	u.Sites.Change(href)
}

func (u *User) String() string {
	return u.Protocol + "=" + u.Subscriber
}

func (u *User) NewUserString() string {
	return fmt.Sprintf("ID: %s"+mapps.Br+
		"Платформа: %s"+mapps.Br+
		"Дата регистрации: %s", mapps.EscapeString(u.Subscriber), mapps.EscapeString(u.Protocol), mapps.EscapeString(time.Unix(u.DateCreated, 0).Format(UserLayout)))
}

func (u *User) FullString(sep string) string {
	return fmt.Sprintf("ID:%s"+sep+
		"Платформа: %s"+sep+
		"Дата регистрации: %s"+sep+
		"Последняя активность: %s"+sep+
		"Уровень допуска: %s", mapps.EscapeString(u.Subscriber), mapps.EscapeString(u.Protocol), mapps.EscapeString(time.Unix(u.DateCreated, 0).Format(UserLayout)), mapps.EscapeString(time.Unix(u.DateLastActivities, 0).Format(UserLayout)), mapps.EscapeString(strconv.Itoa(u.Permission)))
}

type Users struct {
	Mux   sync.RWMutex     `json:"-"`
	Users map[string]*User `json:"users"`
}

func (u *Users) Del(subscriber string) {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	delete(u.Users, subscriber)
}

func (u *Users) User(subscriber string) *User {
	u.Mux.RLock()
	defer u.Mux.RUnlock()
	return u.Users[subscriber]
}

func (u *Users) Add(user *User) {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	if u.Users == nil {
		u.Users = map[string]*User{}
	}

	u.Users[user.Key()] = user
}

func (u *Users) Sample(s func(*User) bool) (result []*User) {
	u.Mux.RLock()
	defer u.Mux.RUnlock()

	for _, user := range u.Users {
		if s(user) {
			result = append(result, user)
		}
	}
	return
}
