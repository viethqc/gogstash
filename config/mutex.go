package config

import "sync"

type mutex struct {
	pause bool
}

var instance *mutex
var once sync.Once

func GetMutexInstance() *mutex {
	once.Do(func() {
		instance = new(mutex)
		instance.pause = false
	})

	return instance
}

func (m *mutex) SetPause(pause bool) {
	m.pause = pause
}

func (m *mutex) GetPause() bool {
	return m.pause
}
