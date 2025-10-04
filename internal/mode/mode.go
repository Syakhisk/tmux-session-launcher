package mode

import "sync"

type Mode string

const (
	ModeAll      Mode = "all"
	ModeSessions Mode = "sessions"
	ModeDir      Mode = "dir"
)

var mu sync.Mutex
var modeCurrent Mode = ModeAll
var modes = []Mode{ModeAll, ModeSessions, ModeDir}

func Get() Mode {
	return modeCurrent
}

func Set(m Mode) {
	modeCurrent = m
}

func Next() Mode {
	mu.Lock()
	defer mu.Unlock()

	for i, m := range modes {
		if m == modeCurrent {
			if i == len(modes)-1 {
				modeCurrent = modes[0]
			} else {
				modeCurrent = modes[i+1]
			}

			break
		}
	}

	return modeCurrent
}

func Prev() Mode {
	mu.Lock()
	defer mu.Unlock()

	for i, m := range modes {
		if m == modeCurrent {
			if i == 0 {
				modeCurrent = modes[len(modes)-1]
			} else {
				modeCurrent = modes[i-1]
			}
			break
		}
	}

	return modeCurrent
}
