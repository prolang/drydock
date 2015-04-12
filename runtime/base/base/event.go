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

// Event represents an auto-reset event used for multi-threaded signalling.
type Event struct {
	// signal is the channel used to indicate when the event is currently signalled.
	signal chan bool

	// signal is an immutable context associated with the Event.  This value is provided at
	// construction and is never mutated by the Event.  It is provided here for application use only.
	data interface{}
}

// NewEvent creates an empty Event with no invalidation or value in the cell.
func NewEvent(data interface{}) *Event {
	// Make a single element channel so that an invalidation can be kept in it and a non-blocking
	// send can used in Put.
	return &Event{
		signal: make(chan bool, 1),
		data:   data,
	}
}

// Signal sends a pulse that will wake up exactly one waiter.
// THREADING: This method is multi-thread safe.
func (e *Event) Signal() {
	select {
	case e.signal <- false:
	default:
	}
}

// Data returns the data element associated with this event.
func (e *Event) Data() interface{} {
	return e.data
}

// Close destroys the event and removes it from any EventSet it is a part of.
func (e *Event) Close() {
	close(e.signal)
}
