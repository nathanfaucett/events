EventEmitter
=====

EventEmitter simple emitter for listening and responding to events

##Example
```
package main

import (
	"github.com/nathanfaucett/event_emitter"
	"fmt"
)

func main() {
	var EE = event_emitter.NewEventEmitter();
	
	EE.On("hello", func(name string) {
		fmt.Println("Hello "+ name)
	})
	
	EE.Emit("hello")
}
```