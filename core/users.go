package core

import (
	"fmt"
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"strconv"
	"sync"
	"time"
)

const (
	UserLayout = "2006/01/02 15:04:05"
)

type User struct {
	mapps.User
	SitesMux           sync.RWMutex `json:"-"`
	Sites              *Set         `json:"sites"`
	Permission         int          `json:"permission"`
	DateCreated        int64        `json:"date_created"`
	DateLastActivities int64        `json:"date_last_activities"`
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

func (u *User) UnSub(href string) {
	u.SitesMux.Lock()
	defer u.SitesMux.Unlock()
	u.Sites.Del(href)
}

func (u *User) ChangeSub(href string) {
	u.SitesMux.Lock()
	defer u.SitesMux.Unlock()
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

func (u *User) FullString() string {
	return fmt.Sprintf("ID:%s"+mapps.Br+
		"Платформа: %s"+mapps.Br+
		"Дата регистрации: %s"+mapps.Br+
		"Последняя активность: %s"+mapps.Br+
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

	u.Users[user.Subscriber] = user
}
