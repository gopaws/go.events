package ticker

import (
	"gopkg.in/ADone/go.events.v1"

	"reflect"
	"sync"
	"time"
)

type needRestart struct{}

func New(emitter events.Emitter) *PeriodicEmitter {
	restart := make(chan needRestart, 1)

	ticker := &PeriodicEmitter{
		Mutex:   new(sync.Mutex),
		Emitter: emitter,
		restart: restart,
		events:  make(map[string]int),
		timers:  []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(restart)}},
	}

	go ticker.run()

	return ticker
}

type PeriodicEmitter struct {
	events.Emitter
	*sync.Mutex
	restart chan needRestart
	events  map[string]int
	timers  []reflect.SelectCase
}

func (emitter *PeriodicEmitter) RegisterEvent(event string, value interface{}, handlers ...events.Listener) {
	ticker.Lock()
	defer ticker.Unlock()

	var timer *time.Ticker
	switch value.(type) {
	case time.Duration:
		timer = time.NewTicker(value.(time.Duration))
	case time.Ticker:
		timer = value.(*time.Ticker)
	default:
		return
	}

	ticker.timers = append(ticker.timers, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timer.C)})

	ticker.events[event] = len(ticker.timers) - 1
	if len(handlers) > 0 {
		ticker.AddEventListener(event, handlers...)
	}

	ticker.restart <- needRestart{}
}

func (ticker *PeriodicEmitter) RemoveEvent(eventName string) {
	ticker.Lock()
	defer ticker.Unlock()

	timerPosition, exists := ticker.events[eventName]
	if !exists {
		return
	}

	delete(ticker.events, eventName)
	ticker.timers = append(ticker.timers[:timerPosition], ticker.timers[timerPosition+1:]...)
	ticker.RemoveEventListeners(eventName)
	ticker.restart <- needRestart{}
}

func (ticker *PeriodicEmitter) run() {
	for {
		selectedIndex, _, ok := reflect.Select(ticker.timers)

		if selectedIndex > 0 {
			for eventName, timerPostion := range ticker.events {
				if ok {
					if timerPostion == selectedIndex {
						ticker.Fire(eventName)
					}
				} else {
					ticker.Lock()
					delete(ticker.events, eventName)
					ticker.timers = append(ticker.timers[:selectedIndex], ticker.timers[selectedIndex+1:]...)
					ticker.Unlock()
				}
			}
		}
	}
}
