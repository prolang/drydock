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

// TurnRunnerSuite is the test suite for Manager
type TurnRunnerSuite struct {
	test.Suite
}

// TestTurnRunnerSuite runs the test suite for Manager
func TestTurnRunnerSuite(t *testing.T) {
	test.RunSuite(t, new(TurnRunnerSuite))
}

func (t *TurnRunnerSuite) Done() async.R {
	return async.Done()
}

func (t *TurnRunnerSuite) DoneWhen() async.R {
	done := async.Done()
	return async.When(done, func(err error) async.R {
		if err != nil {
			t.Errorf("Expected Done() to succeeed.  Got: %v, Want: nil", err)
		}

		return async.Done()
	})
}

func (t *TurnRunnerSuite) NewResultSuccess() async.R {
	r, s := async.NewR()
	async.When(async.Done(), func() {
		s.Complete()
	})
	return r
}

func (t *TurnRunnerSuite) NewResultError() async.R {
	expected := fmt.Errorf("some error")
	r, s := async.NewR()
	async.When(async.Done(), func() {
		s.Fail(expected)
	})
	return async.When(r, func(err error) error {
		if err != expected {
			return fmt.Errorf("Got: %v, Want: %v", err, expected)
		}
		return nil
	})
}

// ResolveCoverage tests Resolver.Resolve().
func (t *TurnRunnerSuite) ResolveCoverage() async.R {
	r, s := async.NewR()
	async.When(async.Done(), func() {
		s.Resolve(nil)
	})

	expected := fmt.Errorf("some error")
	r2, s2 := async.NewR()
	async.When(r, func() {
		s2.Resolve(expected)
	})
	return async.When(r2, func(err error) error {
		if err != expected {
			return fmt.Errorf("Got: %v, Want: %v", err, expected)
		}
		return nil
	})
}

// NewResultForward tests Resolver.Forward().
func (t *TurnRunnerSuite) NewResultForward() async.R {
	expected := fmt.Errorf("some error")
	r, s := async.NewR()
	async.When(async.Done(), func() {
		s.Forward(async.NewError(expected))
	})
	return async.When(r, func(err error) error {
		if err != expected {
			return fmt.Errorf("Got: %v, Want: %v", err, expected)
		}
		return nil
	})
}

// NewErrorfCoverage tests Runner.Errorf().
func (t *TurnRunnerSuite) NewErrorfCoverage() async.R {
	e := async.NewErrorf("some %s %d", "error", 1)
	return async.When(e, func(err error) error {
		if err.Error() != "some error 1" {
			return fmt.Errorf("Got: %v, Want: %v", err, "some error 1")
		}
		return nil
	})
}

// NewResultForwarded tests a When on a forwarded result resolves correctly.
func (t *TurnRunnerSuite) NewResultForwarded() async.R {
	expected := fmt.Errorf("some error")
	r, s := async.NewR()

	// Forward the result synchronously causing the When to be directly enqueued on the error.
	s.Forward(async.NewError(expected))

	return async.When(r, func(err error) error {
		if err != expected {
			return fmt.Errorf("Got: %v, Want: %v", err, expected)
		}
		return nil
	})
}
