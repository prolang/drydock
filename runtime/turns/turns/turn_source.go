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

// This file contains a definition for an I/O Source which is the source of asynchronous
// computations.  Asynchronous computations execute in an environment external to the main turn
// manager of an actor (and so execute independently of the cooperative turn model), but their
// completion can be tracked by and lead to cooperative turn execution through the asynchronous
// result they resolve.

import (
	"sync"

	"github.com/prolang/drydock/runtime/base/base"
	"github.com/prolang/drydock/runtime/turns/async"

	log "github.com/golang/glog"
)

// turnSource represents a source of asynchronous computations whose completions run on a turn
// runner.
type turnSource struct {
	// manager on which turns processing the completion of computation will be executed.
	manager *Manager

	// name is a diagnostic string used to identify the purpose of the turn.
	name string

	// event indicates when there are turns on this source that can be run.
	event *base.Event

	// lock protects list.
	lock sync.Mutex

	// list of turns to be executed on the main runner.
	list *Turn
}

// NewTurnSource creates a new source of I/O computations whose completions run on a turn runner.
func NewTurnSource(runner async.Runner) async.Source {
	manager := runner.(*turnRunner).manager
	t := &turnSource{
		manager: manager,
		name:    "turnSource" + manager.NewID().String(),
		lock:    sync.Mutex{},
		list:    Empty,
	}
	log.Infof("NewTurnSource: %s", t.name)
	t.event = t.manager.registerSource(t)
	return t
}

// Close implements async.Source.Close().
func (t *turnSource) Close() {
	t.event.Close()
}

// New implements async.Source.New().
func (t *turnSource) New(f async.IOFunc) async.R {
	// Allocate a resolver for the caller to use to track the completion of the I/O computation.
	s := newTurnResolver(t.manager)
	r := async.R{async.NewResultT(s)}

	// Pre-allocate a turn from our manager that will execute on the manager later when the I/O
	// computation has completed.
	// THREADING: this turn MUST be allocated here on the manager's thread because manager operations
	// (e.g. NewID() are NOT multi-thread safe).
	var err error
	turn := NewTurn("IOResult"+t.manager.NewID().String(), func() {
		s.Resolve(nil, err)
	})

	go func() {
		// Execute the function on an I/O thread (separate from the turn manager).
		err = f()

		// Once it is finished atomically marshall the result to the I/O source.
		t.lock.Lock()
		t.list = t.list.Append(turn)
		t.lock.Unlock()

		// Signal the source that there is a turn available.
		t.event.Signal()
	}()
	return r
}

// getAllTurns atomically returns all turns that are ready to run (if any).
func (t *turnSource) getAllTurns() /* list */ *Turn {
	t.lock.Lock()
	retval := t.list
	log.V(3).Infof("getAllTurns: %s: %v", t.name, retval)

	t.list = Empty
	t.lock.Unlock()
	return retval
}
