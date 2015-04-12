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
	"reflect"

	"github.com/prolang/drydock/runtime/base/assert"
	"github.com/prolang/drydock/runtime/turns/async"
)

// turnResolver implements the interface Resolver.
type turnResolver struct {
	// manager controls the queue on which turns should be scheduled.
	manager *Manager

	// turns is the queue of turns to be run once resolved in FIFO order.
	turns *Turn

	// outcome is the final value for a completed resolver.  If outcome is NOT of type error then the
	// asynchronous computation completes successfully.  If error then the computation failed.
	outcome interface{}

	// next pointer to another result this result has been forwarded to.
	next *turnResolver
}

// newTurnResolver creates a new unresolved turn-based resolver.
func newTurnResolver(manager *Manager) *turnResolver {
	return &turnResolver{
		manager: manager,
		turns:   Empty,
		outcome: nil,
		next:    nil,
	}
}

// Complete implements Resolver.Complete().
func (s *turnResolver) Complete(value interface{}) {
	_, isError := value.(error)
	assert.True(!isError, "Cannot succeed with an error.")

	s.resolve(value)
}

// Fail implements Resolver.Fail().
func (s *turnResolver) Fail(err error) {
	assert.True(err != nil, "Cannot fail with nil error.")
	s.resolve(err)
}

// Complete implements Resolver.Resolve().
func (s *turnResolver) Resolve(value interface{}, err error) {
	_, isError := value.(error)
	assert.True(!isError, "Cannot succeed with an error.")

	if err != nil {
		assert.True(value == nil, "Provide either value or err but not both.")
		s.resolve(err)
	} else {
		s.resolve(value)
	}
}

// resolve completes the associated Result either successfully or as an error.
func (s *turnResolver) resolve(outcome interface{}) {
	assert.True(!s.isResolved(), "Can't resolve an already resolved result.")

	turns := s.turns
	s.turns, s.outcome = nil, outcome
	s.queueList(turns)
}

// Forward implements Resolver.Forward().
func (s *turnResolver) Forward(n async.ResultT) {
	assert.True(!s.isResolved(), "Cannot forward an already completed result.")

	next := async.InternalUseOnlyGetResolver(n).(*turnResolver)
	assert.True(next.manager == s.manager, "Cannot forward across managers.")

	next = next.getShortest()
	turns := s.turns
	s.turns, s.outcome, s.next = nil, nil, next
	next.queueList(turns)
}

// A set of reflect.Type constants for use in structural validations.
var (
	reflectTypeInterface  = reflect.TypeOf((*interface{})(nil)).Elem()
	reflectTypeError      = reflect.TypeOf((*error)(nil)).Elem()
	reflectTypeAwaitableT = reflect.TypeOf((*async.AwaitableT)(nil)).Elem()
)

// When implements Resolver.WhenT().
func (s *turnResolver) WhenT(in, out, outR reflect.Type, f interface{}) async.ResultT {
	// Validate function type - this is in lieu of static type checking from generics.
	fType := reflect.TypeOf(f)
	assert.True(fType.Kind() == reflect.Func, "f MUST be a WhenFunc")

	// Validate inputs - this is in lieu of static type checking from generics.
	assert.True(fType.NumIn() <= 2, "f MUST take val, err, both or neither")
	takeValue, takeError := true, true
	if numIn := fType.NumIn(); numIn == 2 {
		assert.True(in.AssignableTo(fType.In(0)), "in MUST be assignable to value")
		assert.True(fType.In(1) == reflectTypeError, "f MUST take err")
	} else if numIn == 1 {
		takeError = (fType.In(0) == reflectTypeError)
		takeValue = !takeError
		if takeValue {
			assert.True(in.AssignableTo(fType.In(0)), "in MUST be assignable to value")
		} else {
			assert.True(fType.In(0) == reflectTypeError, "f MUST take err")
		}
	} else if numIn == 0 {
		takeValue, takeError = false, false
	}

	// Validate outputs - this is in lieu of static type checking from generics.
	assert.True(fType.NumOut() <= 2, "f MUST return val, (val, err), or async")
	returnsResult := false
	if numOut := fType.NumOut(); numOut == 2 {
		assert.True(fType.Out(0).AssignableTo(out), "value MUST be assignable to out")
		assert.True(fType.Out(1) == reflectTypeError, "err MUST be error")
	} else if numOut == 1 {
		if fType.Out(0).Implements(reflectTypeAwaitableT) {
			returnsResult = true
			assert.True(fType.Out(0) == outR, "result MUST be ResultT")
		} else if fType.Out(0) == reflectTypeError {
			assert.True(out == reflectTypeInterface, "func() error ONLY allowed on void results")
		} else {
			assert.True(fType.Out(0).AssignableTo(out), "value MUST be assignable to out")
		}
	} else {
		assert.True(out == reflectTypeInterface, "func() ONLY allowed on void results")
	}

	// If this result is already forwarded then just forward the When as well.
	if s.next != nil {
		return s.getShortest().WhenT(in, out, outR, f)
	}

	// Create a new turn that will run once the result is resolved.
	outer := newTurnResolver(s.manager)
	turn := NewTurn("When"+s.manager.NewID().String(), func() {
		// Find the resolved value.
		final := s.getShortest()
		assert.True(final.isResolved(), "When's shouldn't run if the target is not resolved.")

		// Distinguish the error from the value.
		value := final.outcome
		err, isError := final.outcome.(error)
		if isError {
			value = nil
		}

		// If f doesn't take an error but the previous computation failed, then just flow it through to
		// the output by failing immediately.  This allows for null-style error propagation which
		// simplifies handlers since that is a very common case.  f is NEVER called in this case under
		// the assumption that its first line would be: if err != nil { return err }.
		if !takeError {
			if err != nil {
				outer.Fail(err)
				return
			}
		}

		// Convert the arguments into an array of Values
		args := make([]reflect.Value, 0, 2)
		if takeValue {
			if value == nil {
				args = append(args, reflect.New(in).Elem())
			} else {
				args = append(args, reflect.ValueOf(value))
			}
		}
		if takeError {
			if err == nil {
				args = append(args, reflect.New(reflectTypeError).Elem())
			} else {
				args = append(args, reflect.ValueOf(err))
			}
		}

		// Dispatch the result to the WhenFunc.
		retvals := reflect.ValueOf(f).Call(args)

		// Resolve the outer result.
		if returnsResult {
			outer.Forward(retvals[0].Interface().(async.AwaitableT).Base())
			return
		}
		if numRets := len(retvals); numRets == 2 {
			if err, _ = retvals[1].Interface().(error); err != nil {
				outer.resolve(err)
			} else {
				outer.resolve(retvals[0].Interface())
			}
		} else if numRets == 1 {
			outer.resolve(retvals[0].Interface())
		} else {
			outer.resolve(nil)
		}
	})
	s.queue(turn)
	return async.NewResultT(outer)
}

func (s *turnResolver) isResolved() bool {
	return s.turns == nil
}

func (s *turnResolver) getShortest() *turnResolver {
	for s.next != nil {
		s = s.next
	}
	return s
}

func (s *turnResolver) queueList(list *Turn) {
	// Queue all of the pending turns, if there are any, to next.
	var head *Turn
	for !list.IsEmpty() {
		head, list = list.RemoveHead()
		s.queue(head)
	}
}

func (s *turnResolver) queue(turn *Turn) {
	// If the result is already resolved then queue it on the manager for execution, otherwise queue
	// it on the result itself for later.
	if s.isResolved() {
		s.manager.Queue(turn)
	} else {
		s.turns = s.turns.Append(turn)
	}
}
