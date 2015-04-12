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
	"testing"

	"github.com/prolang/drydock/runtime/base/test"

	log "github.com/golang/glog"
)

// RunSuite is the test suite for Turn
type TurnSuite struct {
	test.Suite
}

// TestManagerSuite runs the test suite for Turn
func TestTurnSuite(t *testing.T) {
	test.RunSuite(t, new(TurnSuite))
}

func (t *TurnSuite) Empty() {
	q := Empty

	// Test that the list is empty before any items have been added.
	if !q.IsEmpty() {
		t.Fatalf("Expected empty list.  Got: %v, Want: %v", q, Empty)
	}

	// Test that adding one item makes the list non-empty.
	t1 := NewTurn("t1", func() {})
	q = q.Append(t1)
	log.V(3).Infof("q = %v", q)
	if q.IsEmpty() {
		t.Fatalf("Expected non-empty list.  Got: %v, Want: %v", q, t1)
	}
	if q != t1 {
		t.Fatalf("Expected end of list.  Got: %v, Want: %v", q, t1)
	}
}

func (t *TurnSuite) One() {
	q := Empty
	var head *Turn

	// Test adding one item.
	t1 := NewTurn("t1", func() {})
	q = q.Append(t1)
	log.V(3).Infof("q = %v", q)

	// Test Remove of last item.
	head, q = q.RemoveHead()
	log.V(3).Infof("q = %v", q)
	if !q.IsEmpty() {
		t.Fatalf("Expected empty list.  Got: %v, Want: %v", q, Empty)
	}
	if head != t1 {
		t.Fatalf("Expected pop.  Got: %v, Want: %v", head, t1)
	}
	if head.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", head, nil)
	}

	// Test Unlink of last item.
	q = q.Append(t1)
	log.V(3).Infof("q = %v", q)
	q = q.Unlink(t1)
	log.V(3).Infof("q = %v", q)
	if !q.IsEmpty() {
		t.Fatalf("Expected empty list.  Got: %v, Want: %v", q, Empty)
	}
	if t1.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t1, nil)
	}
}

func (t *TurnSuite) Two() {
	q := Empty
	var head *Turn

	// Test adding two items.
	t1 := NewTurn("t1", func() {})
	t2 := NewTurn("t2", func() {})
	q = q.Append(t2).Append(t1)
	log.V(3).Infof("q = %v", q)
	if q.IsEmpty() {
		t.Fatalf("Expected non-empty list.  Got: %v, Want: %v", q, t1)
	}
	if q != t1 {
		t.Fatalf("Expected end of list.  Got: %v, Want: %v", q, t1)
	}

	// Test Peek returns the correct item when there are two (the head, not the tail).
	head = q.Peek()
	if head != t2 {
		t.Fatalf("Expected head of list.  Got: %v, Want: %v", q, t2)
	}

	// Test removing two items and enusre they come back in the correct order.
	head, q = q.RemoveHead()
	log.V(3).Infof("q = %v", q)
	if head != t2 {
		t.Fatalf("Expected pop.  Got: %v, Want: %v", head, t2)
	}
	head, q = q.RemoveHead()
	log.V(3).Infof("q = %v", q)
	if head != t1 {
		t.Fatalf("Expected pop.  Got: %v, Want: %v", head, t1)
	}
	if !q.IsEmpty() {
		t.Fatalf("Expected empty list.  Got: %v, Want: %v", q, Empty)
	}

	// Test Unlink of two items in {tail, head} order.
	q = q.Append(t2).Append(t1)
	log.V(3).Infof("q = %v", q)
	q = q.Unlink(t1)
	log.V(3).Infof("q = %v", q)
	if q.IsEmpty() {
		t.Fatalf("Expected non-empty list.  Got: %v, Want: %v", q, t2)
	}
	if t1.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t1, nil)
	}
	q = q.Unlink(t2)
	log.V(3).Infof("q = %v", q)
	if !q.IsEmpty() {
		t.Fatalf("Expected empty list.  Got: %v, Want: %v", q, Empty)
	}
	if t2.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t2, nil)
	}

	// Test Unlink of two items in {head, tail} order.
	q = q.Append(t2).Append(t1)
	log.V(3).Infof("q = %v", q)
	q = q.Unlink(t2)
	log.V(3).Infof("q = %v", q)
	if q.IsEmpty() {
		t.Fatalf("Expected non-empty list.  Got: %v, Want: %v", q, t1)
	}
	if t2.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t2, nil)
	}
	q = q.Unlink(t1)
	log.V(3).Infof("q = %v", q)
	if !q.IsEmpty() {
		t.Fatalf("Expected empty list.  Got: %v, Want: %v", q, Empty)
	}
	if t1.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t1, nil)
	}
}

func (t *TurnSuite) Three() {
	q := Empty
	var head *Turn

	// Test adding three items.
	t1 := NewTurn("t1", func() {})
	t2 := NewTurn("t2", func() {})
	t3 := NewTurn("t3", func() {})
	q = q.Append(t3).Append(t2).Append(t1)
	log.V(3).Infof("q = %v", q)
	if q.IsEmpty() {
		t.Fatalf("Expected non-empty list.  Got: %v, Want: %v", q, t1)
	}
	if q != t1 {
		t.Fatalf("Expected end of list.  Got: %v, Want: %v", q, t1)
	}

	// Test Peek returns the correct item when there are three.
	head = q.Peek()
	if head != t3 {
		t.Fatalf("Expected head of list.  Got: %v, Want: %v", head, t3)
	}

	// Test Unlinking an item in the middle of the list.
	q = q.Unlink(t2)
	log.V(3).Infof("q = %v", q)
	head = q.Peek()
	if head != t3 {
		t.Fatalf("Expected head of list.  Got: %v, Want: %v", head, t3)
	}
	head, q = q.RemoveHead()
	log.V(3).Infof("q = %v", q)
	if head != t3 {
		t.Fatalf("Expected pop.  Got: %v, Want: %v", head, t3)
	}
	head, q = q.RemoveHead()
	log.V(3).Infof("q = %v", q)
	if head != t1 {
		t.Fatalf("Expected pop.  Got: %v, Want: %v", head, t1)
	}
	if !q.IsEmpty() {
		t.Fatalf("Expected empty list.  Got: %v, Want: %v", q, Empty)
	}

	// Test Unlink of three items in {tail..head} order.
	q = q.Append(t3).Append(t2).Append(t1)
	log.V(3).Infof("q = %v", q)
	q = q.Unlink(t1)
	log.V(3).Infof("q = %v", q)
	if q.IsEmpty() {
		t.Fatalf("Expected non-empty list.  Got: %v, Want: %v", q, t1)
	}
	if t1.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t1, nil)
	}
	q = q.Unlink(t2)
	log.V(3).Infof("q = %v", q)
	if q.IsEmpty() {
		t.Fatalf("Expected non-empty list.  Got: %v, Want: %v", q, t3)
	}
	if t2.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t2, nil)
	}
	q = q.Unlink(t3)
	log.V(3).Infof("q = %v", q)
	if !q.IsEmpty() {
		t.Fatalf("Expected empty list.  Got: %v, Want: %v", q, Empty)
	}
	if t3.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t3, nil)
	}

	// Test Unlink of three items in {head..tail} order.
	q = q.Append(t3).Append(t2).Append(t1)
	log.V(3).Infof("q = %v", q)
	q = q.Unlink(t3)
	log.V(3).Infof("q = %v", q)
	if q.IsEmpty() {
		t.Fatalf("Expected non-empty list.  Got: %v, Want: %v", q, t1)
	}
	if t3.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t3, nil)
	}
	q = q.Unlink(t2)
	log.V(3).Infof("q = %v", q)
	if q.IsEmpty() {
		t.Fatalf("Expected non-empty list.  Got: %v, Want: %v", q, t1)
	}
	if t2.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t2, nil)
	}
	q = q.Unlink(t1)
	log.V(3).Infof("q = %v", q)
	if !q.IsEmpty() {
		t.Fatalf("Expected empty list.  Got: %v, Want: %v", q, Empty)
	}
	if t1.next != nil {
		t.Fatalf("Expected single item.  Got: %v, Want: %v", t1, nil)
	}

}

func (t *TurnSuite) UnlinkNoop() {
	q := Empty
	t1 := NewTurn("t1", func() {})

	// Test that it is not an error to Unlink a turn that is already Unlinked.  This ensures taht
	// Unlink is idempotent.
	q = q.Unlink(t1)
}

func (t *TurnSuite) String() {
	q := Empty
	if s := q.String(); s == "" {
		t.Errorf("Expected non-empty string representation: Got %s, Want: non-empty", s)
	}

	t1 := NewTurn("t1", func() {})
	if s := t1.Name(); s != "t1" {
		t.Errorf("Expected turn name to round-trip: Got %s, Want: %s", s, "t1")
	}
	if s := t1.String(); s == "" {
		t.Errorf("Expected non-empty string representation: Got %s, Want: non-empty", s)
	}
	q = q.Append(t1)
	if s := q.String(); s == "" {
		t.Errorf("Expected non-empty string representation: Got %s, Want: non-empty", s)
	}

	t2 := NewTurn("t2", func() {})
	q = q.Append(t2)
	if s := q.String(); s == "" {
		t.Errorf("Expected non-empty string representation: Got %s, Want: non-empty", s)
	}

	t3 := NewTurn("t3", func() {})
	q = q.Append(t3)
	if s := q.String(); s == "" {
		t.Errorf("Expected non-empty string representation: Got %s, Want: non-empty", s)
	}

	m := NewManager(NewUniqueIDGenerator())
	m.Queue(NewTurn("t4", func() {}))
	if s := m.String(); s == "" {
		t.Errorf("Expected non-empty string representation: Got %s, Want: non-empty", s)
	}
}
