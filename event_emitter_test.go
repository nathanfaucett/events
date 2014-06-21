package events

import (
	"testing"
)

func TestOn(t *testing.T) {
	event := "test"
	eventEmitter := NewEventEmitter().On(event, func() {})
	
	if len(eventEmitter.events[event]) != 1 {
		t.Error("Failed to add listener to the EventEmitter")
	}
}

func TestOnce(t *testing.T) {
	event := "test"
	emitted := false
	
	eventEmitter := NewEventEmitter().Once(event, func() {
		emitted = true
	}).Emit(event)
	
	if !emitted && len(eventEmitter.events[event]) != 0 {
		t.Error("Failed to add listener to the EventEmitter")
	}
}

func TestOff(t *testing.T) {
	event := "test"
	listener := func() {}
	eventEmitter := NewEventEmitter().On(event, listener).Off(event, listener)
	
	if len(eventEmitter.events[event]) != 0 {
		t.Error("Failed to remove listener from the EventEmitter")
	}
}

func TestRemoveAllListeners(t *testing.T) {
	event := "test"
	eventEmitter := NewEventEmitter().On(event, func() {}).On(event, func() {}).On(event, func() {}).RemoveAllListeners()
	
	if len(eventEmitter.events[event]) != 0 {
		t.Error("Failed to remove all listeners from the EventEmitter")
	}
}

func TestEmit(t *testing.T) {
	event := "test"
	emitted := false
	
	NewEventEmitter().On(event, func() {
		emitted = true
	}).Emit(event)
	
	if !emitted {
		t.Error("Failed to add listener to the EventEmitter")
	}
}