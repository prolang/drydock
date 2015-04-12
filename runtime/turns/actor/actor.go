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

package actor

import (
	"github.com/prolang/drydock/runtime/turns/async"
	"github.com/prolang/drydock/runtime/turns/turns"

	log "github.com/golang/glog"
)

// RunActor starts a contained environment with a turn manager that runs to completion.
// main is the initial turn to be executed and the manager continues execution until main's return
// value is resolved.
func RunActor(root func(async.Runner) async.R) error {
	done := make(chan error, 1)

	go func() {
		manager := turns.NewManager(turns.NewUniqueIDGenerator())
		runner := turns.NewTurnRunner(manager)

		// Allocate a resolver to track the completion of the "main" function.
		r, s := async.NewR(runner)

		// Queue to main routine for execution with the main runner as its runner.
		manager.NewTurn("Main", func() {
			async.When(root(runner), func(err error) {
				log.Infof("Actor Completed with status: %v", err)

				s.Resolve(err)
			})
		})

		err := manager.RunUntil(r)
		done <- err
		close(done)
	}()

	err := <-done
	return err
}
