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

package turns_test

import (
	"fmt"
	"testing"

	"github.com/prolang/drydock/runtime/turns/async"
	"github.com/prolang/drydock/runtime/turns/test"
)

// StringRSuite is the test suite for async.R.
type StringRSuite struct {
	test.Suite
}

// TestStringRSuite runs the test suite for StringRSuite.
func TestStringRSuite(t *testing.T) {
	test.RunSuite(t, new(StringRSuite))
}

// NewCoverage provides coverage for the New function.
func (t *StringRSuite) NewCoverage(runner async.Runner) async.R {
	expected := "a test string"
	r, s := async.NewStringR(runner)
	s.Complete(expected)
	return async.When(r, func(val string, err error) error {
		if err != nil {
			return fmt.Errorf("Expected success.  Got: %v, Want: nil", err)
		}
		if val != expected {
			return fmt.Errorf("Expected success.  Got: %v, Want: %v", val, expected)
		}
		return nil
	})
}

// NewErrorfCoverage provides coverage for the NewErrorf function.
func (t *StringRSuite) NewErrorfCoverage(runner async.Runner) async.R {
	expected := fmt.Errorf("some error")
	r := async.NewStringErrorf(runner, "%v", expected)
	return async.When(r, func(err error) error {
		if err.Error() != expected.Error() {
			return fmt.Errorf("Expected error.  Got: %v, Want: %v", err, expected)
		}
		return nil
	})
}

// WhenCoverage provides coverage for the When function.
func (t *StringRSuite) WhenCoverage(runner async.Runner) async.R {
	expected := "a test string"
	w := async.WhenString(runner.Done(), func() async.StringR {
		r, s := async.NewStringR(runner)
		s.Complete(expected)
		return r
	})

	// Verify that the value was round-tripped.
	return async.When(w, func(val string, err error) error {
		if err != nil {
			return fmt.Errorf("Expected success.  Got: %v, Want: nil", err)
		}
		if val != expected {
			return fmt.Errorf("Expected value.  Got: %v, Want: %v", val, expected)
		}
		return nil
	})
}

// FinallyCoverage provides coverage for the Finally function.
func (t *StringRSuite) FinallyCoverage(runner async.Runner) async.R {
	expected := "a test string"
	didRun := false
	r, s := async.NewStringR(runner)
	s.Complete(expected)
	w := async.FinallyString(r, func() {
		didRun = true
	})

	// Verify that the value was round-tripped.
	return async.When(w, func(val string, err error) error {
		if !didRun {
			return fmt.Errorf("Expected finally to run.  Got: false, Want: true")
		}
		if err != nil {
			return fmt.Errorf("Expected success.  Got: %v, Want: nil", err)
		}
		if val != expected {
			return fmt.Errorf("Expected success.  Got: %v, Want: %v", val, expected)
		}
		return nil
	})
}

// ResolveCoverage provides coverage for the Resolve function.
func (t *StringRSuite) ResolveCoverage(runner async.Runner) async.R {
	expected := "a test string"
	w := async.WhenString(runner.Done(), func() async.StringR {
		r, s := async.NewStringR(runner)
		s.Resolve(expected, nil)
		return r
	})

	// Verify that the value was round-tripped.
	return async.When(w, func(val string, err error) error {
		if err != nil {
			return fmt.Errorf("Expected success.  Got: %v, Want: nil", err)
		}
		if val != expected {
			return fmt.Errorf("Expected value.  Got: %v, Want: %v", val, expected)
		}
		return nil
	})
}

// ForwardCoverage provides coverage for the Forward function.
func (t *StringRSuite) ForwardCoverage(runner async.Runner) async.R {
	expected := "a test string"
	resolved, s0 := async.NewStringR(runner)
	s0.Complete(expected)

	// Instead of resolving the return result, forward it to an already resolved value.
	w := async.WhenString(runner.Done(), func() async.StringR {
		r, s := async.NewStringR(runner)
		s.Forward(resolved)
		return r
	})

	// Verify that the value was round-tripped.
	return async.When(w, func(val string, err error) error {
		if err != nil {
			return fmt.Errorf("Expected success.  Got: %v, Want: nil", err)
		}
		if val != expected {
			return fmt.Errorf("Expected value.  Got: %v, Want: %v", val, expected)
		}
		return nil
	})
}
