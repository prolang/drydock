// Copyright 2015 The Drydock Authors.
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package base

// This file describes an EventSet abstraction.  An EventSet is not itself multithread safe but it
// can be used from a single thread to monitor the activities of others through their events.

import (
	"reflect"

	"github.com/prolang/drydock/runtime/base/assert"
)

// EventSet is a set of events that should be monitored.  It allows zero or more Events to be
// watched for asychronous signalling.
type EventSet struct {
	// The cases using which the EventSet gets its work items
	cases []reflect.SelectCase

	// The events that are being tracked in the set.
	events []*Event
}

// NewEventSet creates a new empty event set.
func NewEventSet() *EventSet {
	// Fill the cases array with the Done context and the internal queueing channel.
	cases := make([]reflect.SelectCase, 1)
	cases[0] = reflect.SelectCase{
		reflect.SelectDefault,
		reflect.ValueOf(nil),
		reflect.ValueOf(nil),
	}
	return &EventSet{
		cases:  cases,
		events: []*Event{nil},
	}
}

// Add registers an event with the set.  Registered events will be returned by Select/Wait when
// they become signaled.
func (w *EventSet) Add(e *Event) {
	// Add the signal channel in the cases array and the event in the events array.
	w.cases = append(w.cases, reflect.SelectCase{
		reflect.SelectRecv,
		reflect.ValueOf(e.signal),
		reflect.ValueOf(nil),
	})
	w.events = append(w.events, e)
}

// Select returns a chosen event or nil if no event is signalled.  Select never blocks.  If no
// event is signalled at the time of the call Select returns nil immediately.  If no events are
// registered then Select return nil immediately.
func (w *EventSet) Select() *Event {
	return w.choose(false)
}

// Wait returns a chosen event or nil if no events are registered.  If no events are signalled at
// the time of the call Wait blocks until an event becomes signalled or all registered events
// become unregistered.
func (w *EventSet) Wait() *Event {
	return w.choose(true)
}

// choose returns the chosen event or nil if no event is ready.
// wait is true if the choose should block for an available event.
// If there are no regiestered then choose always returns immediately with nil.
func (w *EventSet) choose(wait bool) *Event {
	// Configure the default case to block or not depending on wait.
	if wait {
		w.cases[0].Dir = reflect.SelectRecv
	} else {
		w.cases[0].Dir = reflect.SelectDefault
	}

	// Loop around until either:
	// 1.) An event is chosen,
	// 2.) All events have been unregistered,
	// 3.) The default case is selected (if !wait).
	for len(w.cases) > 1 {
		chosen, _, recvOK := reflect.Select(w.cases)

		// If there are no signalled events and we aren't blocking, then return immediately.
		if chosen == 0 {
			return nil
		}

		if recvOK {
			return w.events[chosen]
		}

		// The chosen event needs to be removed.
		w.remove(chosen)
	}

	// No events are registered so return immediately with nil.
	return nil
}

// remove unregisters a closed event from the set.
func (w *EventSet) remove(chosen int) {
	assert.True(chosen != 0, "Never remove the default case")

	// The event has been closed.  Remove it from the set.
	// Put the last element into the relevant index and then remove the last element.
	lastIndex := len(w.events) - 1
	w.events[chosen] = w.events[lastIndex]
	w.events = w.events[:lastIndex]
	w.cases[chosen] = w.cases[lastIndex]
	w.cases = w.cases[:lastIndex]
}
