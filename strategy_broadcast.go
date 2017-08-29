package events

// Broadcast event to all handlers
func Broadcast(event Event, handlers map[Listener]struct{}) {
	for handler := range handlers {
		handler.Handle(event)
	}
}
