package events

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	ErrorInvalidArgument = errors.New("Invalid Argument Kind of Value for listener is not Function")
	// Default number of listeners before warnings
	DefaultMaxListeners = 10
)

type EventEmitter struct {
	*sync.Mutex
	events       map[string][]reflect.Value
	maxListeners int
}

// Creates new EventEmitter
func NewEventEmitter() *EventEmitter {
	this := new(EventEmitter)
	this.Mutex = new(sync.Mutex)
	this.events = make(map[string][]reflect.Value)
	this.maxListeners = DefaultMaxListeners

	return this
}

// attachs listener to this event emitter
func (this *EventEmitter) On(event string, listener interface{}) *EventEmitter {
	this.Lock()
	defer this.Unlock()

	fn := reflect.ValueOf(listener)

	if reflect.Func != fn.Kind() {
		fmt.Println(ErrorInvalidArgument)
		return this
	}
	if this.maxListeners != -1 && this.maxListeners <= len(this.events[event]) {
		fmt.Printf("Warning: event \"%v\" has exceeded the maximum number of listeners of %d\n", event, this.maxListeners)
	}

	this.events[event] = append(this.events[event], fn)

	return this
}

// same as On
func (this *EventEmitter) AddListener(event string, listener interface{}) *EventEmitter {
	return this.On(event, listener)
}

// attachs listener to this event emitter after first fire it removes itself from the events
func (this *EventEmitter) Once(event string, listener interface{}) *EventEmitter {
	fn := reflect.ValueOf(listener)

	if reflect.Func != fn.Kind() {
		fmt.Println(ErrorInvalidArgument)
		return this
	}

	var once func(...interface{})
	once = func(arguments ...interface{}) {
		defer this.Off(event, once)
		var values []reflect.Value

		for i := 0; i < len(arguments); i++ {
			values = append(values, reflect.ValueOf(arguments[i]))
		}

		fn.Call(values)
	}

	return this.On(event, once)
}

// removes listener from this event emitter
func (this *EventEmitter) Off(event string, listener interface{}) *EventEmitter {
	this.Lock()
	defer this.Unlock()

	fn := reflect.ValueOf(listener)

	if reflect.Func != fn.Kind() {
		fmt.Println(ErrorInvalidArgument)
		return this
	}

	if eventList, ok := this.events[event]; ok {
		for i, listener := range eventList {
			if fn == listener {
				this.events[event] = append(this.events[event][:i], this.events[event][i+1:]...)
			}
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

	for event, _ := range this.events {
		this.events[event] = this.events[event][:0]
	}

	return this
}

// emits event and calls all listeners with passed arguments
func (this *EventEmitter) Emit(event string, arguments ...interface{}) *EventEmitter {
	var (
		eventList []reflect.Value
		ok        bool
	)

	this.Lock()
	if eventList, ok = this.events[event]; !ok {
		this.Unlock()
		return this
	}
	this.Unlock()

	var (
		waitGroup sync.WaitGroup
		values    []reflect.Value
	)

	length := len(arguments)
	for i := 0; i < length; i++ {
		values = append(values, reflect.ValueOf(arguments[i]))
	}

	waitGroup.Add(len(eventList))
	for _, fn := range eventList {
		go func(fn reflect.Value) {
			defer waitGroup.Done()
			fn.Call(values)
		}(fn)
	}
	waitGroup.Wait()

	return this
}

// sets maximun number of listeners on events -1 is unlimited
func (this *EventEmitter) SetMaxListeners(max int) *EventEmitter {
	this.maxListeners = max
	return this
}

// returns number of listeners on an event
func (this *EventEmitter) ListenerCount(event string) int {

	return len(this.events[event])
}
