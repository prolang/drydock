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
	"testing"

	"github.com/prolang/drydock/runtime/base/test"
	"github.com/prolang/drydock/runtime/turns/async"
)

// ManagerSuite is the test suite for Manager
type ManagerSuite struct {
	test.Suite
}

// TestManagerSuite runs the test suite for Manager
func TestManagerSuite(t *testing.T) {
	test.RunSuite(t, new(ManagerSuite))
}

func (t *ManagerSuite) QueueOneTurn() {
	m := NewManager(NewUniqueIDGenerator())

	if !m.isIdle() {
		t.Fatalf("Expected empty manager.  Got: %v, Want: %v", m, Empty)
	}

	didRun := false
	f := func() { didRun = true }
	m.Queue(NewTurn("t1", f))
	if m.isIdle() {
		t.Fatalf("Expected non-empty manager.  Got: %v, Want: %v", m, true)
	}

	m.runOneLoop()
	if !m.isIdle() {
		t.Fatalf("Expected empty manager.  Got: %v, Want: %v", m, Empty)
	}
	if !didRun {
		t.Fatalf("Expected turn to have executed.  Got: %v, Want: %v", didRun, true)
	}
}

func (t *ManagerSuite) RunOneTurn() {
	m := NewManager(NewUniqueIDGenerator())

	if !m.isIdle() {
		t.Fatalf("Expected empty manager.  Got: %v, Want: %v", m, Empty)
	}

	didRun := false
	f := func() { didRun = true }
	m.Queue(NewTurn("t1", f))
	if m.isIdle() {
		t.Fatalf("Expected non-empty manager.  Got: %v, Want: %v", m, true)
	}

	ranOne := m.RunOneTurn()
	if !ranOne {
		t.Fatalf("Expected turn to have executed.  Got: %v, Want: %v", ranOne, true)
	}
	if !m.isIdle() {
		t.Fatalf("Expected empty manager.  Got: %v, Want: %v", m, Empty)
	}
	if !didRun {
		t.Fatalf("Expected turn to have executed.  Got: %v, Want: %v", didRun, true)
	}

	ranTwo := m.RunOneTurn()
	if ranTwo {
		t.Fatalf("Expected no turns to have executed.  Got: %v, Want: %v", ranTwo, true)
	}
}

func (t *ManagerSuite) RunUntilSuccess() {
	m := NewManager(NewUniqueIDGenerator())

	// Allocate a resolver to use as the "main" result for RunUntil.
	s := newTurnResolver(m)
	r := async.R{async.NewResultT(s)}

	m.Queue(NewTurn("main", func() {
		s.Complete(nil)
	}))

	// Verify that RunUntil completes when "main" fails.
	err := m.RunUntil(r)
	if err != nil {
		t.Errorf("Expected RunUntil to succeed.  Got: %v, Want: nil", err)
	}
}

func (t *ManagerSuite) RunUntilFailed() {
	m := NewManager(NewUniqueIDGenerator())

	// Allocate a resolver to use as the "main" result for RunUntil.
	s := newTurnResolver(m)
	r := async.R{async.NewResultT(s)}

	expectedError := errors.New("Expected failure")
	m.Queue(NewTurn("main", func() {
		s.Fail(expectedError)
	}))

	// Verify that RunUntil completes when "main" fails.
	err := m.RunUntil(r)
	if err == nil {
		t.Errorf("Expected RunUntil to fail.  Got: %v, Want: %v", err, expectedError)
	}
}

func (t *ManagerSuite) Unlink() {
	m := NewManager(NewUniqueIDGenerator())

	didRun := false
	f := func() { didRun = true }
	t1 := NewTurn("t1", f)
	m.Queue(t1)
	if m.isIdle() {
		t.Fatalf("Expected non-empty manager.  Got: %v, Want: %v", m, true)
	}
	m.Unlink(t1)
	if !m.isIdle() {
		t.Fatalf("Expected empty manager.  Got: %v, Want: %v", m, "empty")
	}
	if didRun {
		t.Fatalf("Expected turn to NOT have executed.  Got: %v, Want: %v", didRun, false)
	}

	didRun = false
	t2 := NewTurn("t2", func() {})
	m.Queue(t2)
	m.Queue(t1)
	m.Unlink(t1)
	m.runOneLoop()
	if didRun {
		t.Fatalf("Expected turn to NOT have executed.  Got: %v, Want: %v", didRun, false)
	}

	didRun = false
	t3 := NewTurn("t3", func() {})
	m.Queue(t2)
	m.Queue(t1)
	m.Queue(t3)
	m.Unlink(t1)
	m.runOneLoop()
	if didRun {
		t.Fatalf("Expected turn to NOT have executed.  Got: %v, Want: %v", didRun, false)
	}
}
