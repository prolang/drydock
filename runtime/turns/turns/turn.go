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
	"fmt"

	"github.com/prolang/drydock/runtime/base/assert"

	log "github.com/golang/glog"
)

// Empty is the empty list of turns.
var Empty = &Turn{}

// Turn represents a bounded synchronous computation to be executed.  A turn's function is executed
// when the turn is run.  A turn may be a member of a list of turns.  A list represents a series of
// computation that are to be performed in FIFO order.  A turn may be in at most one turn list at a
// time.
//
// All operations on turns are O(1) operations with the exception of Unlink() which may be O(n) in
// the length of the list.
type Turn struct {
	// f is the function to execute when running the turn.
	f func()

	// name is a diagnostic string used to identify the purpose of the turn.
	name string

	// next is the next turn if this turn is in a list, otherwise nil.
	next *Turn
}

// NewTurn creates a new single item turn with function f.
func NewTurn(name string, f func()) *Turn {
	return &Turn{
		f:    f,
		name: name,
		next: nil,
	}
}

// Append inserts add at the end of list and returns the new resulting list.
// REQUIRES: add is NOT already in any list.
func (list *Turn) Append(add *Turn) /*newList*/ *Turn {
	assert.True(add != nil, "Can't append null to a turn queue.")
	assert.True(!add.IsEmpty(), "Can't append Empty to a turn queue.")
	assert.True(add.next == nil, "Can't append a turn that is already in a list.")

	log.V(3).Infof("%v: Append %v", list, add)
	if list.IsEmpty() {
		add.next = add // links to itself to complete the circle.
		return add
	}

	head := list.next
	add.next = head
	list.next = add
	return add
}

// Name returns the diagnostic string for the turn.
func (list *Turn) Name() string {
	return list.name
}

// IsEmpty returns true if list is the empty list.
func (list *Turn) IsEmpty() bool {
	return list == Empty
}

// IsList returns true if list is a list, and false if it is a single item.
func (list *Turn) IsList() bool {
	return list.next != nil
}

// Peek returns the turn at the head of the list without removing it from the list.
func (list *Turn) Peek() *Turn {
	return list.next
}

// Unlink removes a turn from a list and returns the new resulting list.
// The turn may appear anywhere in the list including the middle.
// Note: This is potentially an O(n) operation in the length of the list.
// REQUIRES: the item MUST be either a single item or in the list.
func (list *Turn) Unlink(t *Turn) /*newList*/ *Turn {
	// If the item is not in any list, then unlinking is a no-op.
	if !t.IsList() {
		log.V(3).Infof("%v: Unlink no-op %v", list, t)
		return list
	}

	assert.True(!list.IsEmpty(), "Cannot unlink from an empty list.")
	assert.True(list.IsList(), "Cannot unlink from single item.")

	log.V(3).Infof("%v: Unlink %v", list, t)

	// If the list contains only one item then it better be t.
	if list == list.next {
		assert.True(list == t, "Expected item %v to be in list %v.", t, list)

		t.next = nil // unlink the item from the list.
		return Empty
	}

	// If t is the head then just remove it.
	if list.next == t {
		list.next = t.next
		t.next = nil
		return list
	}

	// Find the item (assuming it is in the list).
	before := list.next
	for before.next != t {
		before = before.next
		assert.True(before != list, "Expected item %v to be in list %v.", t, list)
	}
	assert.True(before.next == t, "Expected head %v to be t %v", before.next, t)

	before.next = t.next
	t.next = nil

	// If we removed the list item, the before is the new list, otherwise list hasn't changed.
	if t == list {
		return before
	}
	return list
}

// RemoveHead removes the item from the head of the list and returns both the item removed and the
// new resulting list.
// REQUIRES: The list MUST be non-empty.
func (list *Turn) RemoveHead() ( /*head*/ *Turn /*newList*/, *Turn) {
	assert.True(!list.IsEmpty(), "Cannot remove from empty list.")

	head := list.next

	if list.next == list {
		list = Empty
		head.next = nil
	} else {
		list.next = head.next
		head.next = nil
	}

	log.V(3).Infof("%v: RemoveHead %v", list, head)
	return head, list
}

// Run executes a single turn.
// REQUIRES: the turn be a single turn (i.e. NOT in a list).
func (list *Turn) Run() {
	assert.True(!list.IsEmpty(), "Cannot execute the empty list: %v", list)
	assert.True(list.next == nil, "Cannot execute lists, only single turns: %v.", list)

	log.V(3).Infof("%v: Run", list)

	list.f()
}

// String implements fmt.Stringer
func (list *Turn) String() string {
	// If it is empty then be done.
	if list.IsEmpty() {
		return "<empty>"
	}

	// If the turn is not a list then just print the item.
	if list.next == nil {
		return fmt.Sprintf("<%s, %v, %p>", list.name, list.f, list.next)
	}

	// If a list then print the whole list {head..tail}.
	s := "{ "
	for head := list.next; head != list; head = head.next {
		s += fmt.Sprintf("{%s, %v, %p} ", head.name, head.f, head.next)
	}
	s += fmt.Sprintf("{%s, %v, %p} }", list.name, list.f, list.next)
	return s
}
