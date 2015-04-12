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

import "github.com/prolang/drydock/runtime/turns/async"

// turnRunner is an implementation of async.Runner that uses turns to schedule asynchronous
// computations and completions.
type turnRunner struct {
	manager *Manager
	done    async.R
}

// NewTurnRunner creates a new runner that uses turns to schedule asynchronous
// computations and completions.
func NewTurnRunner(manager *Manager) async.Runner {
	return &turnRunner{
		manager: manager,
		done: async.R{async.NewResultT(&turnResolver{
			manager: manager,
		})},
	}
}

// New implements async.Runner.New().
func (t *turnRunner) New(f async.Func) async.R {
	r, s := async.NewR(t)
	t.manager.NewTurn("New"+t.manager.NewID().String(), func() {
		next := f()
		s.Forward(next)
	})
	return r
}

// NewResult implements async.Runner.NewResultT().
func (t *turnRunner) NewResultT() (async.ResultT, async.ResolverT) {
	s := newTurnResolver(t.manager)
	r := async.NewResultT(s)
	return r, s
}

// Done implements async.Runner.Done().
func (t *turnRunner) Done() async.R {
	return t.done
}
