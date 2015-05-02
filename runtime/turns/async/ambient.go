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

package async

import (
	"runtime"
	"sync"
	"syscall"

	"github.com/prolang/drydock/runtime/base/assert"
)

// This file contains functions for setting and retrieving the ambient (default) runner.  Each actor
// has exactly one ambient runner and in general only the actor framework should ever set one.

type ReleaseFunc func()

// SetAmbientRunner initializes the per-thread storage for an actor context.
// REQUIRES: the caller must call the returned ReleaseFunc to destroy the per-thread context.
func SetAmbientRunner(runner Runner) ReleaseFunc {
	runtime.LockOSThread()
	tid := syscall.Gettid()
	threadContextLock.Lock()
	threadContexts[tid] = threadContext{
		runner: runner,
	}
	threadContextLock.Unlock()
	return func() {
		runtime.UnlockOSThread()
	}
}

// GetCurrentRunner returns the current ambient (default) runner.
func GetCurrentRunner() Runner {
	tid := syscall.Gettid()
	threadContextLock.RLock()
	runner := threadContexts[tid].runner
	threadContextLock.RUnlock()
	assert.True(runner != nil, "GetCurrentRunner can only be called by an actor.")
	return runner
}

// threadContexts contains thread-local storage for use by actors keyed by the OS thread-id of the
// OS thread the actor is running on.  SetAmbientRunner must be called by an actor before any thread
// local storage may be used.
var threadContexts = make(map[int]threadContext)

// threadContextLock protects threadContexts.
var threadContextLock sync.RWMutex

// threadContext contains thread-local storage for use by actors.
type threadContext struct {
	runner Runner
}
