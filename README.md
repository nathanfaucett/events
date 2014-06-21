Events.go
=====

Events.go provides a simple event emitter for listening and responding to events

##Example
```
package main

import (
	"github.com/nathanfaucett/events"
	"fmt"
)

type Parent stuct{
	*events.EventEmitter
}

func NewParent() *Parent {
	this := new(Parent)
	this.EventEmitter = events.NewEventEmitter()
	
	return this
}

func main() {
	parent := NewParent();
	
	parent.On("hello", func(name string) {
		fmt.Println("Hello "+ name)
	})
	
	parent.Emit("hello")
	
	ee := events.NewEventEmitter() //works the same as above
}
```