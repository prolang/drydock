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

// Func is a function that initiates a new computation.  It returns an associated result that is
// resolved when the new computation has completed.
type Func func() R

// Runner abstracts the ability to execute asynchronous computation.  Each computation is
// associated with a result that can be used to monitor its completion.
type Runner interface {
	// New creates a new asynchronous computation and returns its associated result.
	New(Func) R

	// NewResultT creates a new unassociated untyped result and its resolver.
	NewResultT() (ResultT, ResolverT)

	// Done returns an unassociated already successfully completed void result.
	Done() R
}
