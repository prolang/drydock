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

	"github.com/prolang/drydock/runtime/base/base"
	"github.com/prolang/drydock/runtime/base/test"
	"github.com/prolang/drydock/runtime/turns/async"

	log "github.com/golang/glog"
)

// TurnRunnerCoreSuite is the test suite for Manager
type TurnRunnerCoreSuite struct {
	test.Suite
}

// TestTurnRunnerCoreSuite runs the test suite for Manager
func TestTurnRunnerCoreSuite(t *testing.T) {
	test.RunSuite(t, new(TurnRunnerCoreSuite))
}

func (t *TurnRunnerCoreSuite) OneTurn() {
	manager := NewManager(NewUniqueIDGenerator())
	runner := NewTurnRunner(manager)
	tlsRelease := async.SetAmbientRunner(runner)
	defer tlsRelease()

	done := make(chan struct{}, 0)
	async.New(func() async.R {
		close(done)
		return async.Done()
	})
	manager.runOneLoop()
	<-done
}

func (t *TurnRunnerCoreSuite) Done() {
	log.Infof("Running test %s", base.GetMethodName())
	manager := NewManager(NewUniqueIDGenerator())
	runner := NewTurnRunner(manager)
	tlsRelease := async.SetAmbientRunner(runner)
	defer tlsRelease()

	done := async.Done()
	s := async.InternalUseOnlyGetResolver(done.ResultT).(*turnResolver)
	if !s.isResolved() {
		t.Errorf("Expected Done() to be resolved.  Got: %v, Want: resolved", done)
	}
}
