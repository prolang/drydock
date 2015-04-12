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

import (
	"testing"
	"time"

	"github.com/prolang/drydock/runtime/base/test"
)

// EventSetSuite is the test suite for Event and EventSet.
type EventSetSuite struct {
	test.Suite
}

// TestEventSetSuite runs the test suite for Event and EventSet.
func TestEventSetSuite(t *testing.T) {
	test.RunSuite(t, new(EventSetSuite))
}

// SelectUnsignaled verifies that unsignaled events are not chosen by the set.
func (t *EventSetSuite) SelectUnsignaled() {
	e1 := NewEvent(1)
	es := NewEventSet()

	es.Add(e1)
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}
}

// SelectSignaled verifies:
// 1.)  Events are not chosen until signalled.
// 2.)  Events once chosen become unsignalled and aren't chosen again, unless resignalled.
// 3.)  When multiple events are signalled they are all eventually chosen.
func (t *EventSetSuite) SelectSignaled() {
	e1 := NewEvent(1)
	es := NewEventSet()

	// Add one event and verify it is chosen only when it is signalled.
	es.Add(e1)
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}
	e1.Signal()
	if chosen := es.Select(); chosen != e1 {
		t.Errorf("Expected e1.  Got: %v, Want: %v", chosen, e1)
	} else if chosen.Data().(int) != 1 {
		t.Errorf("Expected e1.  Got: %v, Want: %v", chosen.Data(), 1)
	}
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}

	// Add two events and verify only the one signalled is chosen.
	e2 := NewEvent(2)
	es.Add(e2)
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}
	e1.Signal()
	if chosen := es.Select(); chosen != e1 {
		t.Errorf("Expected e1.  Got: %v, Want: %v", chosen, e1)
	}
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}
	e2.Signal()
	if chosen := es.Select(); chosen != e2 {
		t.Errorf("Expected e2.  Got: %v, Want: %v", chosen, e2)
	}
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}

	// Verify that both events are eventually chosen when they are both signalled.
	e1.Signal()
	e2.Signal()
	if chosen1 := es.Select(); chosen1 == nil {
		t.Errorf("Expected either e1 or e2.  Got: nil, Want: non-nil")
	} else if chosen2 := es.Select(); chosen2 == nil {
		t.Errorf("Expected either e1 or e2.  Got: nil, Want: non-nil")
	} else if chosen1 == chosen2 {
		t.Errorf("Expected both e1 or e2.  Got: {%v, %v}, Want: {%v, %v}", chosen1, chosen2, e1, e2)
	}
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}

}

// DoubleSignaled verifies that events that are signalled multiple times between being chosen are
// still chosen only once.
func (t *EventSetSuite) DoubleSignaled() {
	e1 := NewEvent(1)
	es := NewEventSet()

	// Add one event and verify it is chosen only when it is signalled.
	es.Add(e1)
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}
	// Signal the event twice without selecting on it is a no-op.
	e1.Signal()
	e1.Signal()
	if chosen := es.Select(); chosen != e1 {
		t.Errorf("Expected e1.  Got: %v, Want: %v", chosen, e1)
	}
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}
}

// WaitSignaled verifies that Wait blocks until an event is signalled.
func (t *EventSetSuite) WaitSignaled() {
	e1 := NewEvent(1)
	es := NewEventSet()

	doneSleep := false
	es.Add(e1)
	go func() {
		time.Sleep(time.Second)
		doneSleep = true
		e1.Signal()
	}()
	if chosen := es.Wait(); chosen != e1 {
		t.Errorf("Expected e1.  Got: %v, Want: %v", chosen, e1)
	}
	if !doneSleep {
		t.Errorf("Expected Wait to block.  Got: %v, Want: true", doneSleep)
	}
}

// Closed verifies that close removes an event from the set.
func (t *EventSetSuite) Closed() {
	e1 := NewEvent(1)
	es := NewEventSet()

	doneSleep := false
	es.Add(e1)
	go func() {
		time.Sleep(time.Second)
		doneSleep = true
		e1.Close()
	}()
	if chosen := es.Wait(); chosen != nil {
		t.Errorf("Expected nil.  Got: %v, Want: %v", chosen, nil)
	}
	if !doneSleep {
		t.Errorf("Expected Wait to block.  Got: %v, Want: true", doneSleep)
	}

	// Waiting on an empty set always returns immediately.
	if chosen := es.Wait(); chosen != nil {
		t.Errorf("Expected nil.  Got: %v, Want: %v", chosen, nil)
	}
}

// Closed3 veries that events can be removed from the set in any order.
func (t *EventSetSuite) Closed3() {
	e1 := NewEvent(1)
	e2 := NewEvent(2)
	e3 := NewEvent(3)
	es := NewEventSet()

	doneSleep := false
	es.Add(e1)
	es.Add(e2)
	es.Add(e3)
	go func() {
		time.Sleep(time.Second)
		doneSleep = true
		e2.Close()
		e1.Close()
		e3.Close()
	}()
	if chosen := es.Wait(); chosen != nil {
		t.Errorf("Expected nil.  Got: %v, Want: %v", chosen, nil)
	}
	if !doneSleep {
		t.Errorf("Expected Wait to block.  Got: %v, Want: true", doneSleep)
	}

	// Set is not corruptted after removes.
	e4 := NewEvent(4)
	es.Add(e4)
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}
	e4.Signal()
	if chosen := es.Select(); chosen != e4 {
		t.Errorf("Expected e1.  Got: %v, Want: %v", chosen, e4)
	} else if chosen.Data().(int) != 4 {
		t.Errorf("Expected e1.  Got: %v, Want: %v", chosen.Data(), 4)
	}
	if chosen := es.Select(); chosen != nil {
		t.Errorf("Expected no event choosen.  Got: %v, Want: nil", chosen)
	}
}
