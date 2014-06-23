package events

import (
	"errors"
	"fmt"
	"sync"
)

var (
	// Default number of listeners before warnings
	DefaultMaxListeners = 10
)

type EventEmitter struct {
	*sync.Mutex
	events        map[string]*event_list
	max_listeners int
}

// Creates new EventEmitter
func NewEventEmitter() *EventEmitter {
	this := new(EventEmitter)
	this.Mutex = new(sync.Mutex)
	this.events = make(map[string]*event_list)
	this.max_listeners = DefaultMaxListeners

	return this
}

// attachs listener to this event emitter
func (this *EventEmitter) On(event string, listener interface{}) *EventEmitter {
	this.Lock()
	defer this.Unlock()
	var (
		events *event_list
		ok     bool
		length int
	)

	if events, ok = this.events[event]; !ok {
		events = new_event_list()
		this.events[event] = events
	}
	length = events.add(listener)

	if this.max_listeners != -1 && this.max_listeners <= length {
		fmt.Printf("Warning: event \"%v\" has exceeded the maximum number of listeners of %d with %d\n", event, this.max_listeners, length)
	}

	return this
}

// same as On
func (this *EventEmitter) AddListener(event string, listener interface{}) *EventEmitter {
	return this.On(event, listener)
}

// attachs listener to this event emitter after first fire it removes itself from the events
func (this *EventEmitter) Once(event string, listener interface{}) *EventEmitter {
	this.Lock()
	defer this.Unlock()
	var (
		events *event_list
		ok     bool
		length int
	)

	if events, ok = this.events[event]; !ok {
		events = new_event_list()
		this.events[event] = events
	}
	length = events.once(listener)

	if this.max_listeners != -1 && this.max_listeners <= length {
		fmt.Printf("Warning: event \"%v\" has exceeded the maximum number of listeners of %d with %d\n", event, this.max_listeners, length)
	}

	return this
}

// removes listener from this event emitter
func (this *EventEmitter) Off(event string, listener interface{}) *EventEmitter {
	this.Lock()
	defer this.Unlock()

	if events, ok := this.events[event]; ok {
		events.remove(listener)
		if events.length() == 0 {
			delete(this.events, event)
		}
	}

	return this
}

// same as Off
func (this *EventEmitter) RemoveListener(event string, listener interface{}) *EventEmitter {
	return this.Off(event, listener)
}

// removes all events from this
func (this *EventEmitter) RemoveAllListeners() *EventEmitter {
	this.Lock()
	defer this.Unlock()

	for key, events := range this.events {
		events.clear()
		delete(this.events, key)
	}

	return this
}

// emits event and calls all listeners with passed arguments
func (this *EventEmitter) Emit(event string, arguments ...interface{}) *EventEmitter {
	var (
		events *event_list
		ok     bool
	)

	this.Lock()
	if events, ok = this.events[event]; !ok {
		this.Unlock()
		return this
	}
	this.Unlock()
	events.emit(arguments)

	return this
}

// sets maximun number of listeners on events -1 is unlimited
func (this *EventEmitter) SetMaxListeners(max int) *EventEmitter {
	this.max_listeners = max
	return this
}

// returns number of listeners on an event
func (this *EventEmitter) ListenerCount(event string) int {
	this.Lock()
	defer this.Unlock()

	if events, ok := this.events[event]; ok {
		return events.length()
	}

	return 0
}
