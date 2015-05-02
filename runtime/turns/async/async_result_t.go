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

// This file contains internal implementation types for the async package and should not be used
// directly.  Instead see async.R, async.StringR, async.S, async.StringS, etc.

import "reflect"

// ResultT tracks the completion progress of an asynchronous computation.  It is the base type of
// all asynchronous results.  This type is for internal implementation use only.
//
// WARNING: ResultT should NEVER be used directly.  Instead one of the type-safe code generated
// types (e.g. async.R, async.StringR, etc.) should be used.
type ResultT struct {
	s ResolverT
}

// NewResultT allocates a new result for the given resolver.
func NewResultT(s ResolverT) ResultT {
	return ResultT{s}
}

// NewBase creates a new unassociated untyped result and its resolver.
func NewBase() (ResultT, ResolverT) {
	return GetCurrentRunner().NewResultT()
}

// Base provides access to the base ResultT object for subtypes.  This method is intended to be
// inherited by subtypes that implement the AwaitableT interface.
func (r ResultT) Base() ResultT {
	return r
}

// WhenT implements AwaitableT.WhenT().  This method is intended to be inherited by subtypes that
// implement the AwaitableT interface.
func (r ResultT) WhenT(in, out, outR reflect.Type, f interface{}) ResultT {
	return r.s.WhenT(in, out, outR, f)
}

// ResolverT is used to complete an asynchronous computation for which an associated result has been
// allocated.  It is the base type of all asynchronous resolvers.  This type is for internal
// implementation use only.
//
// WARNING: ResolverT should NEVER be used directly.  Instead one of the type-safe code generated
// types (e.g. async.S, async.StringS, etc.) should be used.
type ResolverT interface {
	// Complete completes the associated result successfully with the value provided as the result.
	Complete(value interface{})

	// Fail completes the associated result with an error.
	Fail(err error)

	// Resolve completes the associated result.  If err is nil the result is completed successfully
	// with the value provided as the result, if err is non-nil the result is failed with the error
	// provided and value MUST be nil.
	Resolve(value interface{}, err error)

	// Forward completes the associated result once next is resolved with the same resolution.
	Forward(next ResultT)

	// WhenT schedules a function to be run when the result is completed.
	// See WhenFuncT for the specifications for f.
	WhenT(in, out, outR reflect.Type, f interface{}) ResultT
}

// InternalUseOnlyGetResolver this method is for internal use only and should NEVER be called.
func InternalUseOnlyGetResolver(r ResultT) ResolverT {
	return r.s
}
