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
		events:  make(map[string]*time.Ticker),
		timers:  []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(restart)}},
		mapping: make(map[int]string),
	}

	go ticker.run()

	return ticker
}

type PeriodicEmitter struct {
	events.Emitter
	*sync.Mutex
	restart chan needRestart
	events  map[string]*time.Ticker
	timers  []reflect.SelectCase
	mapping map[int]string
}

func (emitter *PeriodicEmitter) RegisterEvent(event string, value interface{}, handlers ...events.Listener) {
	emitter.Lock()
	defer emitter.Unlock()

	var timer *time.Ticker
	switch value := value.(type) {
	case time.Duration:
		timer = time.NewTicker(value)
	case *time.Ticker:
		timer = value
	default:
		return
	}

	emitter.mapping[len(emitter.timers)] = event
	emitter.timers = append(emitter.timers, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timer.C)})
	emitter.events[event] = timer

	if len(handlers) > 0 {
		emitter.AddEventListener(event, handlers...)
	}

	emitter.restart <- needRestart{}
}

func (emitter *PeriodicEmitter) RemoveEvent(event string) {
	emitter.Lock()
	defer emitter.Unlock()

	_, exists := emitter.events[event]
	if !exists {
		return
	}

	delete(emitter.events, event)
	emitter.refresh()
	emitter.RemoveEventListeners(event)
	emitter.restart <- needRestart{}
}

func (emitter *PeriodicEmitter) refresh() {
	emitter.timers = []reflect.SelectCase{emitter.timers[0]}
	for event, timer := range emitter.events {
		emitter.mapping[len(emitter.timers)] = event
		emitter.timers = append(emitter.timers, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(timer.C),
		})
	}
}

func (emitter *PeriodicEmitter) run() {
	for {
		if index, _, ok := reflect.Select(emitter.timers); index > 0 {
			if event, exists := emitter.mapping[index]; exists {
				if ok {
					emitter.Fire(event)
				} else {
					emitter.Lock()
					delete(emitter.events, event)
					emitter.refresh()
					emitter.Unlock()
				}
			}
		}
	}
}
