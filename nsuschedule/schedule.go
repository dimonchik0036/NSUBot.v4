package nsuschedule

import (
	"encoding/json"
	"errors"
	"github.com/dimonchik0036/Miniapps-pro-SDK"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	emptyLesson = "Свободно"
)

var (
	lessonNumber map[string]int = map[string]int{
		"9:00":  0,
		"10:50": 1,
		"12:40": 2,
		"14:30": 3,
		"16:20": 4,
		"18:10": 5,
		"20:00": 6,
	}

	lessonTime map[int]string = map[int]string{
		0: "09:00-10:35",
		1: "10:50-12:25",
		2: "12:40-14:15",
		3: "14:30-16:05",
		4: "16:20-17:55",
		5: "18:10-19:45",
		6: "20:00-21:35",
	}

	dayString map[int]string = map[int]string{
		0: "Понедельник",
		1: "Вторник",
		2: "Среда",
		3: "Четверг",
		4: "Пятница",
		5: "Суббота",
		6: "Воскресенье",
		/*en:
		- Monday
		- Tuesday
		- Wednesday
		- Thursday
		- Friday
		- Saturday
		- Sunday
		- Empty*/
	}
)

type Schedule struct {
	Mux      sync.RWMutex             `json:"-"`
	Schedule map[string]*ScheduleWeek `json:"schedule"`
}

type ScheduleWeek struct {
	Faculty string       `json:"faculty"`
	Group   string       `json:"group"`
	Lessons [7][]*Lesson `json:"lessons"`
}

func (s *ScheduleWeek) Week() (result [7]string) {
	for day := 0; day < 7; day++ {
		var lessons [7][]*Lesson
		result[day] = dayString[day] + mapps.Br + s.Faculty + mapps.Br + s.Group + mapps.Br + evenString() + mapps.Br
		result[day] = mapps.Bold(result[day])

		for _, l := range s.Lessons[day] {
			number := lessonNumber[l.Time.Start]
			lessons[number] = append(lessons[number], l)
		}

		for i, l := range lessons {
			result[day] += makeLessonString(i, l, false) + mapps.Br
		}
	}

	return
}

func evenString() string {
	_, week := time.Now().ISOWeek()
	if week%2 == 0 {
		return "Чётная (слева)"
	}

	return "Нечётная (справа)"
}

func (s *ScheduleWeek) Day(day int) (result string) {
	var lessons [7][]*Lesson
	now := time.Now().Add(24 * time.Hour * time.Duration(day))
	day = (int(now.Weekday()) + 6) % 7
	result = dayString[day] + mapps.Br + s.Faculty + mapps.Br + s.Group + mapps.Br

	_, week := now.ISOWeek()
	var even bool
	if week%2 == 0 {
		even = true
		result += "Чётная (слева)" + mapps.Br
	} else {
		result += "Нечётная (справа)" + mapps.Br
	}
	result = mapps.Bold(result)

	for _, l := range s.Lessons[day] {
		number := lessonNumber[l.Time.Start]
		lessons[number] = append(lessons[number], l)
	}

	for i, l := range lessons {
		result += makeLessonString(i, l, even) + mapps.Br
	}

	return
}

func makeLessonString(number int, lessons []*Lesson, even bool) string {
	if len(lessons) == 0 {
		return mapps.Bold(strconv.Itoa(number+1)) + " " + mapps.Bold(lessonTime[number]) + " " + mapps.Italic(emptyLesson)
	}

	if len(lessons) == 1 && lessons[0].Date.Week == 0 {
		return mapps.Bold(strconv.Itoa(number+1)) + " " + mapps.Bold(lessonTime[number]) + " " + lessons[0].String()
	}

	var evenLesson *Lesson
	var oddLesson *Lesson
	for _, l := range lessons {
		if l.Date.Week == 1 {
			oddLesson = l
		} else {
			evenLesson = l
		}
	}

	return mapps.Bold(strconv.Itoa(number+1)) + " " + mapps.Bold(lessonTime[number]) + " " + helpLesson(oddLesson, even) + " / " + helpLesson(evenLesson, even)
}

func helpLesson(lesson *Lesson, even bool) string {
	if lesson == nil {
		return "~"
	}

	if even && lesson.Date.Week == 2 {
		return mapps.Bold(lesson.String())
	}

	if !even && lesson.Date.Week == 1 {
		return mapps.Bold(lesson.String())
	}

	return lesson.String()
}

func makeScheduleWeek(faculty string, group *Group) (week *ScheduleWeek) {
	week = new(ScheduleWeek)
	week.Faculty = faculty
	week.Group = group.Name
	week.Lessons = group.GetLessons()
	return
}

func (s *Schedule) GetGroup(group string) (*ScheduleWeek, bool) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	str, ok := s.Schedule[group]
	return str, ok
}

func (s *Schedule) GetDay(group string, day int) (string, bool) {
	g, ok := s.GetGroup(group)
	if !ok {
		return "", ok
	}

	return g.Day(day), true
}

const (
	scheduleUrl = "http://table.nsu.ru/api"
)

func getSchedule() (University, error) {
	code, body, err := fasthttp.Get(nil, scheduleUrl)
	if err != nil {
		return University{}, err
	}

	if code != fasthttp.StatusOK {
		return University{}, errors.New("Bad code: " + strconv.Itoa(code))
	}

	var u University
	if err := json.Unmarshal(body, &u); err != nil {
		return University{}, err
	}

	return u, nil
}

func NewSchedule() Schedule {
	u, err := getSchedule()
	if err != nil {
		return Schedule{}
	}

	return Schedule{Schedule: u.MakeSchedule()}
}

func (s *Schedule) Update() {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	u, err := getSchedule()
	if err != nil {
		log.Print("Bad schedule update: ", err)
		return
	}

	s.Schedule = u.MakeSchedule()
}

type University struct {
	Name      string       `json:"name"`
	Abbr      string       `json:"abbr"`
	Faculties []*Faculties `json:"faculties"`
}

func (u *University) MakeSchedule() map[string]*ScheduleWeek {
	s := make(map[string]*ScheduleWeek)
	for _, f := range u.Faculties {
		for _, g := range f.Groups {
			s[g.Name] = makeScheduleWeek(f.Name, g)
		}
	}

	return s
}

type Faculties struct {
	Name   string   `json:"name"`
	Groups []*Group `json:"groups"`
}

type Group struct {
	Name    string    `json:"name"`
	Lessons []*Lesson `json:"lessons"`
}

func (g *Group) GetLessons() (lessons [7][]*Lesson) {
	for _, l := range g.Lessons {
		if l.Date.Weekday > 7 || l.Date.Weekday < 1 {
			continue
		}

		lessons[l.Date.Weekday-1] = append(lessons[l.Date.Weekday-1], l)
	}
	return
}

type Lesson struct {
	Subject string `json:"subject"`
	Type    string `json:"type"`
	Time    struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"time"`
	Date struct {
		Start   string `json:"start"`
		End     string `json:"end"`
		Weekday int    `json:"weekday"`
		Week    int    `json:"week"`
	} `json:"date"`
	Audiences []struct {
		Name string `json:"name"`
	} `json:"audiences"`
	Teachers []struct {
		Name string `json:"name"`
	} `json:"teachers"`
}

func (l *Lesson) String() string {
	return l.Type + ", " + l.Subject + ", " + help(l.Audiences) + ", " + help(l.Teachers)
}

func help(s []struct {
	Name string `json:"name"`
}) string {
	var res []string
	for _, n := range s {
		res = append(res, n.Name)
	}
	return strings.Join(res, ", ")
}
