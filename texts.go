package main

import (
	"github.com/dimonchik0036/Miniapps-wrapper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
)

const (
	LangRu = "ru"
	LandEn = "en"

	LangDefault = LangRu
)

var TextsForAdmin Texts
var TextsForUsers Texts

const (
	FilenameTextsUsers = "texts_for_users"
	FilenameTextsAdmin = "texts_for_admin"
)

var FakeText *Text = &Text{
	Title:      "Здесь должен был быть текст, но его нет",
	Button:     "Кнопка, которой нет",
	BackButton: "Путь в никуда",
	Navigation: `<navigation><link pageId="` + StrPageMenuMain + `">Дорога в никуда</link> </navigation>`,
	Optional:   []string{"Reload: "},
}

type Texts struct {
	Mux   sync.RWMutex             `yaml:"-"`
	Texts map[string]*Localization `yaml:"texts"`
}

func (t *Texts) Add(key string, lang string, text *Text) {
	t.Mux.Lock()
	defer t.Mux.Unlock()
	if t.Texts == nil {
		t.Texts = map[string]*Localization{}
	}

	loc, ok := t.Texts[key]
	if ok {
		loc.Add(lang, text)
		return
	}
	loc = new(Localization)
	loc.Add(lang, text)

	t.Texts[key] = loc
}

func (t *Texts) Get(key string, lang string) *Text {
	t.Mux.RLock()
	defer t.Mux.RUnlock()
	if t.Texts == nil {
		t.Texts = map[string]*Localization{}
	}

	loc, ok := t.Texts[key]
	if !ok {
		Admin.SendMessage(key + " not found!")
		return FakeText
	}

	return loc.Get(lang)
}

type Localization struct {
	Mux  sync.RWMutex     `yaml:"-"`
	Lang map[string]*Text `yaml:",inline"`
}

func (t *Localization) Add(lang string, text *Text) {
	t.Mux.Lock()
	defer t.Mux.Unlock()

	if t.Lang == nil {
		t.Lang = map[string]*Text{}
	}
	t.Lang[lang] = text
}

func (t *Localization) Get(lang string) *Text {
	t.Mux.RLock()
	defer t.Mux.RUnlock()
	if t.Lang == nil {
		return nil
	}

	text := t.Lang[lang]
	if text == nil {
		text = t.Lang[LangDefault]
	}
	if text == nil {
		return FakeText
	}

	return text
}

type Text struct {
	Title      string   `yaml:"title,omitempty"`
	Body       string   `yaml:"body,omitempty"`
	Button     string   `yaml:"button,omitempty"`
	BackButton string   `yaml:"back_button,omitempty"`
	Navigation string   `yaml:"navigation,omitempty"`
	Optional   []string `yaml:"optional,omitempty"`
}

func (t *Text) DoPage(args string) string {
	return mapps.Page(args,
		mapps.Title("",
			t.Title,
		),
		t.Body,
		t.Navigation,
	)
}

func (t *Text) GetOptional(index int) string {
	if index >= 0 && index < len(t.Optional) {
		return t.Optional[index]
	}

	Admin.SendMessageTelegram("Button: " + t.Button + " optional out of range " + strconv.Itoa(index))
	return ""
}

func loadAllTexts(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename + ".yaml")
	if err != nil {
		Admin.SendMessageTelegram("Ошибка открытия текстов: " + err.Error())
		log.Print(err)
		return
	}
	TextsForUsers.Mux.Lock()
	defer TextsForUsers.Mux.Unlock()

	if err := yaml.Unmarshal(data, v); err != nil {
		Admin.SendMessageTelegram("Ошибка расшифровки: " + err.Error())
		log.Print(err)
		return
	}

	f, err := os.OpenFile(filename+"_new.yaml", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0700))
	if err != nil {
		Admin.SendMessageTelegram("Ошибка шифровки: " + err.Error())
		log.Print(err)
		return
	}
	defer f.Close()
	data, err = yaml.Marshal(v)
	if err != nil {
		Admin.SendMessageTelegram("Ошибка шифровки: " + err.Error())
		log.Print(err)

		return
	}

	f.Write(data)
}
