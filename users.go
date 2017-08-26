package main

import (
	. "github.com/dimonchik0036/Miniapps-pro-SDK"
	"github.com/dimonchik0036/NSUBot/core"
	"time"
)

var GlobalSubscribers *core.Users

func CheckNewSubscriber(request *Request) *core.User {
	user := GlobalSubscribers.User(request.Subscriber)
	if user == nil {
		return newUser(request.User())
	}
	return user
}

func newUser(user User) (subscriber *core.User) {
	subscriber = new(core.User)
	subscriber.Queue.Lock()
	defer subscriber.Queue.Unlock()
	subscriber.User = user
	GlobalSubscribers.Add(subscriber)

	subscriber.Permission = core.PermissionUser
	t := time.Now().Unix()
	subscriber.DateCreated = t
	subscriber.DateLastActivities = t

	go Admin.SendMessageBlock(Div("", subscriber.NewUserString()))
	return
}

func refreshSubscriber(subscriber *core.User) *core.User {
	subscriber.DateLastActivities = time.Now().Unix()
	return subscriber
}
