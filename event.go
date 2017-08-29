package events

// EventOption for event
type EventOption struct {
	apply func(*Event)
}

// WithContext sets event metadata
func WithContext(context Map) EventOption {
	return EventOption{func(event *Event) {
		for key, value := range context {
			event.Context[key] = value
		}
	}}
}

// New create new event with provided name and options
func New(data interface{}, options ...EventOption) Event {
	var event Event

	switch value := data.(type) {
	case string:
		event = Event{value, Map{}}
	case Event:
		event = value
	}

	for _, option := range options {
		option.apply(&event)
	}

	return event
}

// Event
type Event struct {
	Key     string
	Context Map
}

func (event *Event) String() string {
	return event.Key
}
