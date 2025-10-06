package mode

import "sync"

type Mode string

const (
	ModeAll       Mode = "all"
	ModeSession   Mode = "session"
	ModeDirectory Mode = "directory"
)

var mu sync.Mutex
var modeCurrent = ModeAll
var Modes = []Mode{ModeAll, ModeSession, ModeDirectory}

func Get() Mode {
	return modeCurrent
}

func Set(m Mode) {
	modeCurrent = m
}

func Next() Mode {
	mu.Lock()
	defer mu.Unlock()

	for i, m := range Modes {
		if m == modeCurrent {
			if i == len(Modes)-1 {
				modeCurrent = Modes[0]
			} else {
				modeCurrent = Modes[i+1]
			}

			break
		}
	}

	return modeCurrent
}

func Prev() Mode {
	mu.Lock()
	defer mu.Unlock()

	for i, m := range Modes {
		if m == modeCurrent {
			if i == 0 {
				modeCurrent = Modes[len(Modes)-1]
			} else {
				modeCurrent = Modes[i-1]
			}
			break
		}
	}

	return modeCurrent
}

func (m Mode) String() string {
	return string(m)
}
