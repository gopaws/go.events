package events

// ParallelBroadcast calls event handlers in separate goroutines
func ParallelBroadcast(event Event, handlers map[Listener]struct{}) {
	for handler := range handlers {
		go handler.Handle(event)
	}
}
