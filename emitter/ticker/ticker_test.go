package ticker

import (
	"gopkg.in/ADone/go.events.v1"

	"testing"
	"time"
)

type TestDispatcher struct {
	Count  int
	Target bool
}

func (handler *TestDispatcher) AddSubscribers(t ...events.Listener) {
}

func (handler *TestDispatcher) Dispatch(event events.Event) {
	handler.Count++
	handler.Target = true
}

func TestNewPeriodicEmitter(t *testing.T) {
	emitter := New()

	if emitter.started != false {
		t.Log("Fail ticker create - corrupted ran flag")
		t.Fail()
	}
	if len(emitter.Timers) != 0 {
		t.Log("Fail ticker create - no timers storage")
		t.Fail()
	}
	if len(emitter.Dispatchers) != 0 {
		t.Log("Fail ticker create - no dispatcher storage")
		t.Fail()
	}

}

func TestStart(t *testing.T) {
	emitter := New()

	if emitter.started != false {
		t.Fail()
	}
	emitter.Start()
	if emitter.started != true {
		t.Log("Fail ticker start")
		t.Fail()
	}

}

func TestRegisterEvent(t *testing.T) {
	emitter := New()
	dispatcher := new(TestDispatcher)

	if dispatcher.Target != false {
		t.Fail()
	}

	emitter.RegisterEvent("test", 1*time.Millisecond)
	emitter.Dispatchers["test"] = dispatcher
	emitter.Start()
	time.Sleep(3 * time.Millisecond)

	if dispatcher.Target != true {
		t.Log("fail event fire")
		t.Fail()
	}

	if dispatcher.Count < 2 {
		t.Log("fail event fire - not enough events")
		t.Fail()
	}
}

func TestAddEventListener(t *testing.T) {
	emitter := New()

	if len(emitter.Dispatchers) != 0 {
		t.Fail()
	}
	listener := events.Callback(func(event events.Event) {})

	emitter.AddEventListener("test_1", listener)
	emitter.AddEventListener("test_2", listener)

	if len(emitter.Dispatchers) != 2 {
		t.Log("Fail add listeners")
		t.Fail()
	}
}
