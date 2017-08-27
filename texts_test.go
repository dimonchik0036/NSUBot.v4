package main

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestTexts(t *testing.T) {
	texts := Texts{}
	texts.Add("menu_main", "ru", FakeText)
	texts.Add("menu_main", "en", FakeText)
	texts.Add("option", "ru", FakeText)

	data, err := yaml.Marshal(texts)
	if err != nil {
		t.Fatal(err)
	}
	println(string(data))
}
