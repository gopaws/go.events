package events_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"../events"
)

var _ = Describe("Event", func() {
	It("should create event object", func() {
		Expect(events.New("test")).To(BeEquivalentTo(events.Event{Key: "test", Context: events.Map{}}))
	})

	It("should subscribe callback to event", func() {
		emitter := events.NewEmitter()
		emitter.On("test", events.Callback(func(event events.Event) {}))
		// Expect(events.New("test")).To(BeEquivalentTo(events.Event{Key: "test", Context: events.Map{}}))
	})
})
