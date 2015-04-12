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

package turns

import (
	"errors"
	"fmt"

	"github.com/prolang/drydock/runtime/base/assert"
	"github.com/prolang/drydock/runtime/base/base"
	"github.com/prolang/drydock/runtime/turns/async"
)

// Manager is a queue of turns that can be executed either one at a time or all together.  New
// turns can be added to the managers queue at any time.
type Manager struct {
	// sources is the set of I/O sources from which asynchronous turns may arrive.
	sources *base.EventSet

	// turns is the main queue of turns to be run by this manager in FIFO order.
	turns *Turn

	// idgen generates new unique ids.
	idgen *UniqueIDGenerator
}

// NewManager creates a new turn manager.
func NewManager(idgen *UniqueIDGenerator) *Manager {
	return &Manager{
		sources: base.NewEventSet(),
		turns:   Empty,
		idgen:   idgen,
	}
}

// NewID generates a new ID.  ID's are never reused.
func (m *Manager) NewID() UniqueID {
	return m.idgen.NewID()
}

// String implements fmt.Stringer
func (m *Manager) String() string {
	return fmt.Sprintf("%v", m.turns)
}

// Queue adds an existing turn to the managers queue.
// REQUIRES: t is NOT already in any turn list or queue.
func (m *Manager) Queue(t *Turn) {
	m.turns = m.turns.Append(t)
}

// NewTurn creates a new turn that will call f() when it is executed.  Adds the turn to the
// managers queue.
func (m *Manager) NewTurn(name string, f func()) {
	t := NewTurn(name, f)
	m.turns = m.turns.Append(t)
}

// Unlink removes a turn from the manager.
// The turn may appear anywhere in the managers queue including the middle.
// Note: This is potentially an O(n) operation in the length of the queue.
// REQUIRES: the turn MUST be in the queue.
func (m *Manager) Unlink(t *Turn) {
	m.turns = m.turns.Unlink(t)
}

// errSuccess is a sentinel value used to identify a successful exit.
var errSuccess = errors.New("main exited successfully")

// RunUntil runs turns in the manager until main becomes resolved.  If main fails, then its error
// is returned, otherwise returns nil.
func (m *Manager) RunUntil(main async.R) error {
	var mainExited error
	m.NewTurn("RunUntil", func() {
		async.When(main, func(err error) error {
			if err != nil {
				mainExited = err
			} else {
				mainExited = errSuccess
			}
			return err
		})
	})

	// Loop around until the main result has been resolved.  Block efficiently on I/O (so that we
	// don't spin on select) if we run out of local work to do.
	for mainExited == nil {
		// Flush the main queue.
		m.runOneLoop()

		// If there is no work to do then block on I/O.
		if m.turns.IsEmpty() && (mainExited == nil) {
			e := m.sources.Wait()
			assert.True(m.turns.IsEmpty(), "Only blocked on I/O if there was no work to do.")
			assert.True(mainExited == nil, "Only blocked on I/O if the program not exited.")
			assert.True(e != nil, "LIVE LOCK: no progress can be made because there are no I/O "+
				"sources,	no local turns, and the program has not yet exited.")
			e.Signal() // force Select in next loop to see this source again.
		}
	}

	if mainExited != errSuccess {
		return mainExited
	}
	return nil
}

// RunOneTurn runs a single turn from the main queue if any exist.  Return true if a turn was run.
func (m *Manager) RunOneTurn() bool {
	var t *Turn
	if !m.turns.IsEmpty() {
		t, m.turns = m.turns.RemoveHead()
		t.Run()
		return true
	}
	return false
}

// runOneLoop runs a single iteration of the turn loop which includes both executing turns on the
// main queue at the time of the call and checking for new turns from asynchronous sources.
func (m *Manager) runOneLoop() {
	// Check for async I/O turns and append them to the main queue before snapshotting.
	for e := m.sources.Select(); e != nil; e = m.sources.Select() {
		var head *Turn
		list := e.Data().(*turnSource).getAllTurns()
		for !list.IsEmpty() {
			head, list = list.RemoveHead()
			m.turns = m.turns.Append(head)
		}
	}

	// Snapshot the main queue.  Executing these turns may enqueue more turns on the main queue but
	// they won't be part of this loop (but will be part of the next loop).
	list := m.turns
	m.turns = Empty

	// Run a full cycle of the main queue snapshot.
	var t *Turn
	for !list.IsEmpty() {
		t, list = list.RemoveHead()
		t.Run()
	}
}

// isIdle returns true if the manager's main queue has no turns to run, otherwise false.
func (m *Manager) isIdle() bool {
	return m.turns.IsEmpty()
}

// registerSource adds an asynchronous source of turns to the set tracked by this manager.
// Returns an event the source can use to signal when turns are available.  The returned event
// should be closed to unregister the source.
func (m *Manager) registerSource(source *turnSource) *base.Event {
	// Allocate the event HERE in manager instead of in turnSource because manager does the down-cast
	// on event.Data() and MUST ensure it will be the type expected.  If the event were passed in by
	// the caller then the caller might incorrectly set event.Data() to some other value or type.
	event := base.NewEvent(source)

	// Add the event to the set of source to track.  Closing the event will automatically unregister
	// the source during the main turn loop's Select call.
	m.sources.Add(event)
	return event
}
