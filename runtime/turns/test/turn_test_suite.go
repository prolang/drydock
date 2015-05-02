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

package test

// This file contains the definition of test.Suite which extends base/test.Suite and runs a series
// of turn-based tests as a suite.

import (
	"reflect"
	"testing"

	"github.com/prolang/drydock/runtime/base/assert"
	"github.com/prolang/drydock/runtime/base/test"
	"github.com/prolang/drydock/runtime/turns/actor"
	"github.com/prolang/drydock/runtime/turns/async"
)

// Suite is the base type for turn-based test package.  All test suites using the test package
// should embed this value.
type Suite struct {
	test.Suite
}

// RunSuite runs all test methods within a test suite and reports their outcome to stdout.
func RunSuite(t *testing.T, suite interface{}) {
	test.RunSuiteCustom(t, suite, filterTurnTests, dispatchTurnTests)
}

// filterTurnTests filters tests to only those that are either synchronous or asynchronous
// turn-based test methods.
func filterTurnTests(m reflect.Method) bool {
	// Synchronous test methods are also allowed.
	if m.Type.NumIn() == 1 && m.Type.NumOut() == 0 {
		return true
	}

	// Asynchronous test method are allowed.
	if m.Type.NumIn() != 1 || m.Type.NumOut() != 1 {
		return false
	}
	if m.Type.Out(0) != asyncResultType {
		return false
	}

	return true
}

// A set of reflect.Type constants for use in structural validations.
var (
	asyncTestFuncType = reflect.TypeOf((*async.Func)(nil)).Elem()
	asyncResultType   = reflect.TypeOf((*async.R)(nil)).Elem()
)

// dispatchTurnTests dispatches both synchronous and turn-based asynchronous tests.
func dispatchTurnTests(s *test.Suite, v reflect.Value, f reflect.Value) {
	assert.True(f.Kind() == reflect.Func, "Test function MUST be a function")

	// If it is a synchronous test method then just run it.
	if f.Type().NumOut() == 0 {
		inputs := []reflect.Value{v}
		f.Call(inputs)
		return
	}

	// If an asynchronous test method then run it within a new Actor and wait for the Actor to
	// complete.
	fn := reflect.MakeFunc(asyncTestFuncType, func(args []reflect.Value) []reflect.Value {
		// Prepend the receiver argument to the function before dispatching.
		inputs := []reflect.Value{v}
		inputs = append(inputs, args...)
		return f.Call(inputs)
	})

	err := actor.RunActor(fn.Interface().(async.Func))
	if err != nil {
		s.Errorf("Expected test case retval to succeed.  Got: %q, Want: nil", err)
	}
}
