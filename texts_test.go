package main

import (
	"github.com/dimonchik0036/Miniapps-wrapper"
	"testing"
)

func TestTexts(t *testing.T) {
	println(mapps.Page("",
		mapps.Input("submit", StrPageFeedback, "."),
		mapps.Navigation(mapps.FormatAttr("id", "submit"),
			mapps.Link("", StrPageFeedback, "submit")),
		mapps.Navigation("",
			mapps.Link("",
				StrPageMenuOption, "В меню")),
	))
}
