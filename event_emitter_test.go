package events

import (
	"errors"
	"strconv"
	"testing"
)

func BenchmarkEvents(b *testing.B) {
	eventEmitter := NewEventEmitter()
	eventEmitter.On("test", func(err error, num int, name string) {
		for i := 0; i < 100; i++ {
			if err == nil {
				name = name + strconv.Itoa(num)
			} else {
				name = err.Error() + strconv.Itoa(num)
			}
		}
	})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eventEmitter.Emit("test", nil, 10, "fun")
		eventEmitter.Emit("test", errors.New("Not Fun"), 10, "fun")
	}
}

func TestOn(t *testing.T) {
	event := "test"
	eventEmitter := NewEventEmitter().On(event, func() {})

	if eventEmitter.ListenerCount(event) != 1 {
		t.Error("Failed to add listener to the EventEmitter")
	}
}

func TestOnce(t *testing.T) {
	event := "test"
	emitted := false

	eventEmitter := NewEventEmitter().Once(event, func() {
		emitted = true
	}).Emit(event)

	if !emitted && eventEmitter.ListenerCount(event) != 0 {
		t.Error("Failed to add listener to the EventEmitter")
	}
}

func TestOff(t *testing.T) {
	event := "test"
	listener := func() {}
	eventEmitter := NewEventEmitter().On(event, listener).Off(event, listener)

	if eventEmitter.ListenerCount(event) != 0 {
		t.Error("Failed to remove listener from the EventEmitter")
	}
}

func TestRemoveAllListeners(t *testing.T) {
	event := "test"
	eventEmitter := NewEventEmitter().On(event, func() {}).On(event, func() {}).On(event, func() {}).RemoveAllListeners()

	if eventEmitter.ListenerCount(event) != 0 {
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
