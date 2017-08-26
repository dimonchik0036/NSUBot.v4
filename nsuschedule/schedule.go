package nsuschedule

import "sync"

type Schedule struct {
	Mux      sync.RWMutex      `json:"-"`
	Schedule map[string]string `json:"schedule"`
}

func GetAllSchedule() (Schedule, error) {
	return NewSchedule(), nil
}

func NewSchedule() Schedule {
	return Schedule{
		Schedule: make(map[string]string),
	}
}

func (s *Schedule) Update() {
	s.Mux.Lock()
	defer s.Mux.Unlock()
}
