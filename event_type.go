package events

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
)

type event_list struct {
	listeners []reflect.Value
	arguments []reflect.Type
	initted   bool
}

func new_event_list() *event_list {
	this := new(event_list)
	return this
}

func (this *event_list) init(function interface{}) {
	typeof := reflect.TypeOf(function)
	if typeof.Kind() != reflect.Func {
		fmt.Println(ErrorInvalidArgument)
		os.Exit(1)
	}
	var arguments []reflect.Type
	length := typeof.NumIn()
	for i := 0; i < length; i++ {
		arguments = append(arguments, typeof.In(i))
	}

	this.arguments = arguments
	this.initted = true
}

func (this *event_list) add(function interface{}) int {
	if !this.initted {
		this.init(function)
	}
	fn := reflect.ValueOf(function)

	if fn.Kind() != reflect.Func {
		fmt.Println(ErrorInvalidArgument)
		os.Exit(1)
	}
	length := fn.Type().NumIn()
	if length != len(this.arguments) {
		fmt.Println(ErrorInvalidArgument)
		os.Exit(1)
	}
	this.listeners = append(this.listeners, fn)

	return len(this.listeners)
}

func (this *event_list) once(function interface{}) int {
	if !this.initted {
		this.init(function)
	}
	fn := reflect.ValueOf(function)

	if fn.Kind() != reflect.Func {
		fmt.Println(ErrorInvalidArgument)
		os.Exit(1)
	}
	length := fn.Type().NumIn()
	if length != len(this.arguments) {
		fmt.Println(ErrorInvalidArgument)
		os.Exit(1)
	}

	var once func(...interface{})
	once = func(arguments ...interface{}) {
		defer this.remove(once)
		var values []reflect.Value

		for i, argument := range arguments {
			if argument == nil {
				values = append(values, reflect.Zero(this.arguments[i]))
			} else {
				values = append(values, reflect.ValueOf(argument))
			}
		}

		fn.Call(values)
	}
	this.listeners = append(this.listeners, reflect.ValueOf(once))

	return len(this.listeners)
}

func (this *event_list) remove(function interface{}) {
	fn := reflect.ValueOf(function)

	for i, listener := range this.listeners {
		if fn == listener {
			this.listeners = append(this.listeners[:i], this.listeners[i+1:]...)
		}
	}
}

func (this *event_list) emit(arguments []interface{}) {
	length := len(this.listeners)
	if length == 0 {
		return
	}
	var values []reflect.Value

	for i, argument := range arguments {
		if argument == nil {
			values = append(values, reflect.Zero(this.arguments[i]))
		} else {
			values = append(values, reflect.ValueOf(argument))
		}
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(length)

	for _, listener := range this.listeners {
		go func(fn reflect.Value) {
			defer func() {
				waitGroup.Done()
				if r := recover(); nil != r {
					err := errors.New(fmt.Sprintf("%v", r))
					fmt.Println(err)
				}
			}()
			fn.Call(values)
		}(listener)
	}
	waitGroup.Wait()
}

func (this *event_list) clear() {
	this.listeners = this.listeners[:0]
}

func (this *event_list) length() int {
	return len(this.listeners)
}
