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

// This file contains the definition of test.Suite which runs a series of tests as a suite with the
// possibility of common storage across the tests.  The test.Suite framework also provides for test
// name encapsulation to avoid test name collision across suites in the same package.

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/prolang/drydock/runtime/base/assert"

	log "github.com/golang/glog"
)

// Suite is the base type for the test package.  All test suites using the test package should
// embed this value.
type Suite struct {
	t         *testing.T
	suiteName string
	failed    bool
	msg       string
}

// Infof implements log.Infof for a Suite.
func (t *Suite) Infof(format string, args ...interface{}) {
	log.InfoDepth(1, fmt.Sprintf(format, args...))
}

// Errorf implements log.Errorf for a Suite and fails the current test and suite.
func (t *Suite) Errorf(format string, args ...interface{}) {
	t.failed = true
	t.msg = t.getFileLine() + " " + fmt.Sprintf(format, args...)
	log.ErrorDepth(1, fmt.Sprintf(format, args...))
	t.t.Fail()
}

// Fatalf implements log.Fatalf for a Suite and fails the current test and suite.
func (t *Suite) Fatalf(format string, args ...interface{}) {
	t.failed = true
	t.msg = t.getFileLine() + " " + fmt.Sprintf(format, args...)
	log.ErrorDepth(1, fmt.Sprintf(format, args...))
	t.t.FailNow()
}

// RunSuite runs all test methods within a test suite and reports their outcome to stdout.
func RunSuite(t *testing.T, suite interface{}) {
	RunSuiteCustom(t, suite, filterSyncTests, dispatchSyncTests)
}

// DispatchFunc is a function that provides custom test dispatching.
// s is the base Suite from which logging and failure can be recorded.
// v is the concrete suite instance passed to RunSuite.
// f is the test function being executed.
type DispatchFunc func(s *Suite, v reflect.Value, f reflect.Value)

// FilterFunc is a function for selecting which test methods to execute.
// f is the test function to be filtered.
// returns true if the test should be executed, false if the test should be skipped.
type FilterFunc func(m reflect.Method) bool

// RunSuiteCustom uses dispatch to run all test methods within a test suite identified by filter
// and reports their outcome to stdout.
func RunSuiteCustom(t *testing.T, suite interface{}, filter FilterFunc, dispatch DispatchFunc) {
	v := reflect.ValueOf(suite)

	// Walk upward the type embedding tree until we find the most-root inclusion of test.Suite.
	base := v.Elem().FieldByName("Suite")
	for base.Type() != reflect.TypeOf((*Suite)(nil)).Elem() {
		// If this is a suite that embeds another suite look for its suite's test.Suite
		base = base.FieldByName("Suite")
	}
	s := base.Addr().Interface().(*Suite)
	s.t = t

	// Discover the list of methods defined in the base type Suite so these won't be executed as
	// test methods.
	skip := make(map[string]bool)
	for i := 0; i < reflect.TypeOf(s).NumMethod(); i++ {
		m := reflect.TypeOf(s).Method(i)
		skip[m.Name] = true
	}

	// Fetch the pointer and non-pointer versions of the type metadata for enumeration purposes.
	ptrType := v.Type()
	suiteType := v.Elem().Type()

	log.Infof("Running Suite: %v, %s, %s", ptrType, suiteType.Name(), suiteType.PkgPath())
	s.suiteName = suiteType.Name()

	// Iterate through all of the test methods and execute them in order.
	for i := 0; i < ptrType.NumMethod(); i++ {
		m := ptrType.Method(i)

		// Skip methods "inherited" from Suite.
		if skip[m.Name] {
			continue
		}

		// Skip non-public methods.
		if m.Name[:1] != strings.ToUpper(m.Name[:1]) {
			continue
		}

		// Skip filtered methods
		if !filter(m) {
			continue
		}

		// Execute a test method and then print out its outcome.
		startTime := time.Now()
		fmt.Printf("=== RUN %s.%s\n", s.suiteName, m.Name)
		log.Infof("=== RUN %s.%s\n", s.suiteName, m.Name)

		s.failed = false
		dispatch(s, v, m.Func)
		outcome := "PASS"
		if s.failed {
			outcome = "FAIL"
		}
		endTime := time.Now()
		duration := endTime.Sub(startTime).Seconds()

		log.Infof("--- %s: %s.%s (%.2fs)\n", outcome, s.suiteName, m.Name, duration)
		fmt.Printf("--- %s: %s.%s (%.2fs)\n", outcome, s.suiteName, m.Name, duration)
		if s.failed {
			fmt.Println(s.msg)
		}
	}
}

// filterSyncTests selects only those test methods that take no arguments and have no return value.
func filterSyncTests(m reflect.Method) bool {
	// Skip methods with the wrong signature.
	if m.Type.NumIn() != 1 || m.Type.NumOut() != 0 {
		return false
	}

	return true
}

// dispatchSyncTests dispatches a synchronous test method that takes no arguments and has no return
// value.  Dispatch waits for the test method to complete.
func dispatchSyncTests(s *Suite, v reflect.Value, f reflect.Value) {
	assert.True(f.Kind() == reflect.Func, "Test function MUST be a function")

	inputs := []reflect.Value{v}
	f.Call(inputs)
}

// getFileLine returns the caller's file and line number in the format file:line.
func (t *Suite) getFileLine() string {
	_, file, line, ok := runtime.Caller(2)
	assert.True(ok, "The caller MUST always be available for user code")

	return fmt.Sprintf("%s:%d", file, line)
}
