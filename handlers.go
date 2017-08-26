package main

import (
	"fmt"
	. "github.com/dimonchik0036/Miniapps-pro-SDK"
	"github.com/dimonchik0036/NSUBot/news"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

func NewsHandler(subscribers []string, news []news.News, title string) {
	for _, s := range subscribers {
		subscriber := GlobalSubscribers.User(s)
		if subscriber == nil {
			log.Printf("%s %s", subscriber.String(), "WTF?! nil pointer in newshandler: "+title)
			continue
		}

		for _, n := range news {
			if err := subscriber.SendMessageBlock(Div("", EscapeString(title)+Br+EscapeString(n.URL)+Br+Br+stringCheck(n.Title)+stringCheck(n.Decryption)+EscapeString(time.Unix(n.Date, 0).Format("02.01.2006")))); err != nil {
				log.Printf("%s %s", subscriber.String(), err.Error())
			}
		}
	}
}

func MainHandler(request *Request) {
	subscriber := CheckNewSubscriber(request)
	subscriber.Queue.Lock()
	defer subscriber.Queue.Unlock()
	SystemLoger.Print(request.AllFields())
	log.Print(subscriber.String() + " " + request.String())
	refreshSubscriber(subscriber)

	switch string(request.Page) {
	default:
		fmt.Fprint(request.Ctx,
			Page("",
				Div("", EscapeString(time.Now().Format(time.RFC3339))),
				Navigation("",
					Link("",
						"One",
						"_1"),
					Link("",
						"Two",
						"_2"),
				),
			),
		)
	}
}

func stringCheck(s string) string {
	if s == "" {
		return ""
	} else {
		return EscapeString(s) + Br
	}
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	req, err := Decode(string(ctx.RequestURI()))
	if err != nil {
		log.Print(err, " ", string(ctx.RequestURI()))
		fmt.Fprint(ctx, "404 Not Found")
		return
	}

	req.Ctx = ctx
	MainHandler(&req)
}

func HandlerStart() {
	Admin.SendMessage("Начинаю!")
	fasthttp.ListenAndServe(Port, fastHTTPHandler)
}
