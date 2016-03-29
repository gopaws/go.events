go.events
===========

`go.events` is a small [Observer](https://en.wikipedia.org/wiki/Observer_pattern) implemetation for golang

[![GoDoc](https://godoc.org/github.com/ADone/go.events?status.svg)](https://godoc.org/github.com/ADone/go.events)
[![Join the chat at https://gitter.im/ADone/go.events](https://badges.gitter.im/ADone/go.events.svg)](https://gitter.im/ADone/go.events?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Import
------

`go.events` available through [gopkg.in](http://labix.org/gopkg.in) interface:
```go
import "gopkg.in/ADone/go.events.v1"
```
or directly from github:
```go
import "github.com/ADone/go.events"
```

Usage
-----

### Event

Creating standalone event object:
```go
event := events.New("eventName")
event.Meta["key"] = value
```
For `Meta` field see [go.meta](https://github.com/ADone/go.meta)

### Emitter

Package `emiter` implements `events.Emitter` interface
```go
import (
	"gopkg.in/ADone/go.events.v1"
	"gopkg.in/ADone/go.events.v1/emitter"
)
```

#### Create

Emitter could combined with other structs via common `events.Emitter` interface:
```go
type Object struct {
	events.Emitter
}

object := Object{emitter.New()}
```
> it's preferable usage example,
> it simplify test cases of base structs

Emitter could be created with specific dispatch strategy:
```go
import "gopkg.in/ADone/go.events.v1/dispatcher"
```

``` go
emitter.New(dispatcher.BroadcastFactory)
emitter.New(dispatcher.ParallelBroadcastFactory)
```

#### Emmit event

Emit concrete event object:

```go
em := emitter.New()
em.Fire(events.New("event"))
```

Emit event with label & params:
```go
em.Fire("event")
// or with event params
em.Fire("event", meta.Map{"key": "value"})
// or with plain map
em.Fire("event", map[string]interface{}{"key": "value"})
````
> Be carefully with concurrent access to `event.Meta`

#### Subscribe for event

Emitter supports only `events.Listener` interface for subscription, but it can be extended by embedded types:

* channels
```go
channel := make(chan events.Event)

object.AddEventListener("event", events.Stream(channel))
```
* handlers
```go
type Handler struct {}

func (Handler) Handle (events.Event) {
	// handle events
}

object.AddEventListener("event", Handler{})
// or
object.On("event", Handler{}, Handler{}).On("anotherEvent", Handler{})
```
* functions
```go
object.AddEventListener("event", events.Callback(func(event events.Event){
	// handle event
}))
```

### Ticker
Package `ticker` adds support of periodic events on top of events.Emitter

```go
import (
	"gopkg.in/ADone/go.events.v1/emitter/ticker"
	"gopkg.in/ADone/go.events.v1/emitter"
	"time"
)
```

```go
tick := ticker.New(emitter.New())
tick.RegisterEvent("periodicEvent1", 5*time.Second)
// or
tick.RegisterEvent("periodicEvent2", time.NewTicker(5*time.Second))
// or directly with handlers
tick.RegisterEvent("periodicEvent3", 5*time.Second, Handler{})
```
