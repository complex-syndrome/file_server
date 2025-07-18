package helper

import (
	"sync"
)

type FanOut struct {
	_subscribers []chan string
	_lock        sync.Mutex
}

func (f* FanOut) Subscribe() <-chan string {
	nchan := make(chan string, 100)
	f._lock.Lock()
	defer f._lock.Unlock()

	f._subscribers = append(f._subscribers, nchan)
	return nchan
}

func (f* FanOut) Publish(msg string) {
	f._lock.Lock()
	defer f._lock.Unlock()

	for _, subscriber := range f._subscribers {
		select {
		case subscriber <- msg:
		default:
		}	
	}
}