package events

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
)

var (
	// Default number of listeners before warnings
	DefaultMaxListeners  = 10
	ErrorInvalidArgument = errors.New("Invalid Argument listener is not a Function")
)

type event_listener struct {
	function  reflect.Value
	arguments []reflect.Type
}

func new_event_listener(function interface{}) *event_listener {
	fn := reflect.ValueOf(function)

	if fn.Kind() != reflect.Func {
		fmt.Println(ErrorInvalidArgument)
		os.Exit(1)
	}

	var arguments []reflect.Type

	typeof := fn.Type()
	length := typeof.NumIn()
	for i := 0; i < length; i++ {
		arguments = append(arguments, typeof.In(i))
	}

	this := new(event_listener)
	this.function = fn
	this.arguments = arguments

	return this
}

func (this *event_listener) values(arguments []interface{}) []reflect.Value {
	var values []reflect.Value

	for i, argument := range arguments {
		if argument == nil {
			values = append(values, reflect.Zero(this.arguments[i]))
		} else {
			values = append(values, reflect.ValueOf(argument))
		}
	}

	return values
}

type EventEmitter struct {
	*sync.Mutex
	events       map[string][]*event_listener
	maxListeners int
}

// Creates new EventEmitter
func NewEventEmitter() *EventEmitter {
	this := new(EventEmitter)
	this.Mutex = new(sync.Mutex)
	this.events = make(map[string][]*event_listener)
	this.maxListeners = DefaultMaxListeners

	return this
}

// attachs listener to this event emitter
func (this *EventEmitter) On(event string, listener interface{}) *EventEmitter {
	this.Lock()
	defer this.Unlock()

	fn := new_event_listener(listener)
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
	this.Lock()
	defer this.Unlock()

	fn := new_event_listener(listener)

	if this.maxListeners != -1 && this.maxListeners <= len(this.events[event]) {
		fmt.Printf("Warning: event \"%v\" has exceeded the maximum number of listeners of %d\n", event, this.maxListeners)
	}

	var warp func(...interface{})
	warp = func(arguments ...interface{}) {
		defer this.Off(event, warp)
		fn.function.Call(fn.values(arguments))
	}

	once := new_event_listener(warp)
	once.arguments = fn.arguments

	this.events[event] = append(this.events[event], once)

	return this
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
			if fn == listener.function {
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
		eventList []*event_listener
		ok        bool
	)

	this.Lock()
	if eventList, ok = this.events[event]; !ok {
		this.Unlock()
		return this
	}
	this.Unlock()

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(eventList))
	for _, listener := range eventList {
		values := listener.values(arguments)

		go func(fn reflect.Value) {
			defer func() {
				if r := recover(); nil != r {
					err := errors.New(fmt.Sprintf("%v", r))
					fmt.Println(err)
				}
			}()
			defer waitGroup.Done()
			fn.Call(values)
		}(listener.function)
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
