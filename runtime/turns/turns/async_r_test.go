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

// RSuite is the test suite for async.R.
type RSuite struct {
	test.Suite
}

// TestRSuite runs the test suite for RSuite.
func TestRSuite(t *testing.T) {
	test.RunSuite(t, new(RSuite))
}

// NewCoverage provides coverage for the New function.
func (t *RSuite) NewCoverage(runner async.Runner) async.R {
	r, s := async.NewR(runner)
	s.Complete()
	return async.When(r, func(val interface{}, err error) async.R {
		if err != nil {
			return async.NewErrorf(runner, "Expected successful completion.  Got: %v, Want: nil", err)
		}
		return runner.Done()
	})
}

// NewErrorfCoverage provides coverage for the NewErrorf function.
func (t *RSuite) NewErrorfCoverage(runner async.Runner) async.R {
	expected := fmt.Errorf("some error")
	r := async.NewErrorf(runner, "%v", expected)
	return async.When(r, func(val interface{}, err error) async.R {
		if err.Error() != expected.Error() {
			return async.NewErrorf(runner, "Expected error.  Got: %v, Want: %v", err, expected)
		}
		return runner.Done()
	})
}

// WhenCoverage provides coverage for the When function.
func (t *RSuite) WhenCoverage(runner async.Runner) async.R {
	w := async.When(runner.Done(), func(val interface{}, err error) async.R {
		r, s := async.NewR(runner)
		s.Complete()
		return r
	})

	// Verify full signature.
	expected := interface{}(nil)
	w1 := async.When(w, func(val interface{}, err error) async.R {
		if err != nil {
			return async.NewErrorf(runner, "Expected success w1.  Got: %v, Want: nil", err)
		}
		if val != expected {
			return async.NewErrorf(runner, "Expected value w1.  Got: %v, Want: %v", val, expected)
		}
		return w
	})

	// Verify just error.
	w2 := async.When(w1, func(err error) async.R {
		if err != nil {
			return async.NewErrorf(runner, "Expected success w2.  Got: %v, Want: nil", err)
		}
		return w1
	})

	// Verify just value.
	w3 := async.When(w2, func(val interface{}) async.R {
		if val != expected {
			return async.NewErrorf(runner, "Expected value w3.  Got: %v, Want: %v", val, expected)
		}
		return w2
	})

	// Verify neither value nor error.
	didRun := false
	w4 := async.When(w3, func() async.R {
		didRun = true
		return w3
	})
	return async.When(w4, func() async.R {
		if !didRun {
			return async.NewErrorf(runner, "Expected executed w4.  Got: %v, Want: true", didRun)
		}
		return w4
	})
}

// WhenErrorCoverage provides coverage for the When function.
func (t *RSuite) WhenErrorCoverage(runner async.Runner) async.R {
	expected := fmt.Errorf("some error")
	w := async.When(runner.Done(), func() async.R {
		r, s := async.NewR(runner)
		s.Fail(expected)
		return r
	})

	w1 := async.When(w, func(_ interface{}, err error) async.R {
		if err != expected {
			return async.NewErrorf(runner, "Expected executed w.  Got: %v, Want: %v", err, expected)
		}
		return w
	})

	w2 := async.When(w1, func(interface{}) async.R {
		return async.NewErrorf(runner, "Expected not executed w1")
	})
	w3 := async.When(w2, func(_ interface{}, err error) async.R {
		if err != expected {
			return async.NewErrorf(runner, "Expected executed w.  Got: %v, Want: %v", err, expected)
		}
		return w2
	})

	w4 := async.When(w3, func() async.R {
		return async.NewErrorf(runner, "Expected not executed w3")
	})
	return async.When(w4, func(_ interface{}, err error) async.R {
		if err != expected {
			return async.NewErrorf(runner, "Expected executed w.  Got: %v, Want: %v", err, expected)
		}
		return runner.Done()
	})
}

// WhenReturnCoverage provides coverage for the When function.
func (t *RSuite) WhenReturnCoverage(runner async.Runner) async.R {
	w := async.When(runner.Done(), func() async.R {
		r, s := async.NewR(runner)
		s.Complete()
		return r
	})

	w1 := async.When(w, func() (interface{}, error) {
		return nil, nil
	})

	w2 := async.When(w1, func() error {
		return nil
	})

	w3 := async.When(w2, func() interface{} {
		return nil
	})

	w4 := async.When(w3, func() {
	})

	return w4
}

// WhenErrorReturnCoverage provides coverage for the When function.
func (t *RSuite) WhenErrorReturnCoverage(runner async.Runner) async.R {
	expected := fmt.Errorf("some error")
	w := async.When(runner.Done(), func() async.R {
		r, s := async.NewR(runner)
		s.Fail(expected)
		return r
	})

	w1 := async.When(w, func(err error) (interface{}, error) {
		if err != expected {
			return nil, fmt.Errorf("Expected executed w.  Got: %v, Want: %v", err, expected)
		}
		return nil, err
	})

	w2 := async.When(w1, func(err error) error {
		if err != expected {
			return fmt.Errorf("Expected executed w1.  Got: %v, Want: %v", err, expected)
		}
		return err
	})

	return async.When(w2, func(err error) async.R {
		if err != expected {
			return async.NewErrorf(runner, "Expected executed w2.  Got: %v, Want: %v", err, expected)
		}
		return runner.Done()
	})
}

// FinallyCoverage provides coverage for the Finally function.
func (t *RSuite) FinallyCoverage(runner async.Runner) async.R {
	expected := fmt.Errorf("some error")
	didRun := false
	r, s := async.NewR(runner)
	s.Fail(expected)
	w := async.Finally(r, func() {
		didRun = true
	})

	// Verify that the value was round-tripped.
	return async.When(w, func(val interface{}, err error) async.R {
		if !didRun {
			return async.NewErrorf(runner, "Expected finally to run.  Got: false, Want: true")
		}
		if err != expected {
			return async.NewErrorf(runner, "Expected error.  Got: %v, Want: %v", err, expected)
		}
		if val != nil {
			return async.NewErrorf(runner, "Expected error.  Got: %v, Want: nil", val)
		}
		return runner.Done()
	})
}
