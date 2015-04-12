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
	"testing"

	"github.com/prolang/drydock/runtime/turns/async"
	"github.com/prolang/drydock/runtime/turns/test"
	"github.com/prolang/drydock/runtime/turns/turns"
)

// TurnSourceSuite is the test suite for TurnSource.
type TurnSourceSuite struct {
	test.Suite
}

// TestTurnSourceSuite runs the test suite for TurnSource.
func TestTurnSourceSuite(t *testing.T) {
	test.RunSuite(t, new(TurnSourceSuite))
}

// AsyncTurn verifies that a parallel turn executed by a Source does run, and that its resolution is
// reflected correctly on the main turn queue.
func (t *TurnSourceSuite) AsyncTurn(runner async.Runner) async.R {
	src := turns.NewTurnSource(runner)

	syncPoint1 := make(chan struct{}, 0)
	syncPoint2 := make(chan struct{}, 0)

	// Create a truly parallel turn.
	r := src.New(func() error {
		t.Infof("Beginning async turn.")

		// produce sync point 1.
		close(syncPoint1)

		// wait for sync point 2.
		<-syncPoint2

		return nil
	})

	// Create a coop-turn to provide the sync point.
	runner.New(func() async.R {
		t.Infof("Beginning coop turn that provides sync point.")
		// Wait for both parallel and coop turns to have started.
		<-syncPoint1

		// Completed the parallel turn.
		close(syncPoint2)
		return runner.Done()
	})

	// Cleanup the source.
	return async.Finally(r, func() {
		// Destroy the TurnSource since we don't need it anymore.
		src.Close()
	})
}
