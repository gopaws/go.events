package events

type Map map[string]interface{}

// DispatchStrategy defines strategy of delivery event to handlers
type DispatchStrategy func(Event, map[Listener]struct{})

// Listener defines event handler interface
type Listener interface {
	Handle(Event)
}

// Stream implements Listener interface on channel
type Stream chan Event

// Handle Listener
func (stream Stream) Handle(event Event) {
	stream <- event
}

// Callback implements Listener interface on function
func Callback(function func(Event)) Listener {
	return callback{&function}
}

type callback struct {
	function *func(Event)
}

// Handle Listener
func (callback callback) Handle(event Event) {
	(*callback.function)(event)
}
