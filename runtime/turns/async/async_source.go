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

// IOFunc is a function that initiates a new computation on an I/O thread.  It returns error that
// is used to resolve an associated R.
type IOFunc func() error

// Source represents a source of truly parallel computations from outside a coop Runner.  Source can
// be used to bridge between cooperative concurrent execution and truly parallel computations.
type Source interface {
	// New starts a new I/O computation and returns an associated result.
	New(f IOFunc) R

	// Close destroys the I/O source.
	//
	// WARNING: Any outstanding I/O's will NOT complete their results and will be orphaned.  The
	// caller is responsible for waiting for any pending I/O to complete *before* calling Close().
	Close()
}
