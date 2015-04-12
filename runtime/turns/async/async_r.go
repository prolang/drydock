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

// The types async.R and async.S are the most basic types in the async package.  These types
// represent a asynchronous computation with no return value (i.e. void-returning).  This represents
// the most basic structure of an asynchronous computation of the form:
//
//   func() error
//
// All asynchronous computations execute and then complete.  When completing (either successfully or
// as an error) they resolve their associated result.  Asynchronous computation with no return value
// (other than a possible error) are called void-returning and their completion is indicated by the
// resolution of an async.R.
//
// The most basic way to start an asynchronous computation is directly from the current Actor's
// Runner:
//
//   r := runner.New(func() async.R {
//          // do something
//          return runner.Done()
//		    })
//
// In this example the func() will be executed asynchronous and r will become resolved when both
// func() and all of its child computations have completed.
//
// Though the above is the most basic formulation, the most *common* way to start an asynchronous
// computation is by expressing a dependency through a When operation:
//
//   r2 := async.When(r, func() {
//           // do something else once "r" is complete.
//         })
//
// In this When example the func() will be executed when the computation tracked by r has completed.
// Thus r2 becomes a dependent of r and begins only when r is complete.  Unless the error value from
// the dependency is explicitly requested, a dependent computation will be automatically aborted if
// its dependency fails.  Alternatively a dependent computation can examine the outcome of its
// dependency and provide logic on both the success and failure paths:
//
//   r2 := async.When(r, func(err error) {
//           if err != nil {
//	           // do something different if "r" failed.
//	           return doRecovery()
//           }
//           // do something else once "r" is complete.
//         })
//

import (
	"fmt"
	"reflect"
)

// R tracks the completion progress of an asynchronous computation.
type R struct {
	ResultT
}

// Type implements AwaitableT.Type().
func (R) Type() reflect.Type {
	return reflect.TypeOf((*interface{})(nil)).Elem()
}

// NewR allocates a new result.
func NewR(runner Runner) (R, S) {
	r, s := runner.NewResultT()
	return R{r}, S{s}
}

// NewError returns an unassociated already failed result.
func NewError(runner Runner, err error) R {
	r, s := NewR(runner)
	s.Fail(err)
	return r
}

// NewErrorf returns an unassociated already failed result.
func NewErrorf(runner Runner, format string, a ...interface{}) R {
	return NewError(runner, fmt.Errorf(format, a...))
}

// When implements AwaitableT.WhenT().
func When(r AwaitableT, f interface{}) R {
	reflectTypeInterface := reflect.TypeOf((*interface{})(nil)).Elem()
	reflectTypeR := reflect.TypeOf((*R)(nil)).Elem()
	return R{r.WhenT(r.Type(), reflectTypeInterface, reflectTypeR, f)}
}

// Finally schedules a function to be run regardless of how an R resolves.  The outcome of a finally
// operation is always the same as the original R.
func Finally(r R, f func()) R {
	return When(r, func(interface{}, error) R {
		f()
		return r
	})
}

// S is used to complete an asynchronous computation for which an associated R has been allocated.
type S struct {
	s ResolverT
}

// Complete implements ResolverT.Complete() for void.
func (r S) Complete() {
	r.s.Complete(nil)
}

// Fail implements ResolverT.Fail() for void.
func (r S) Fail(err error) {
	r.s.Fail(err)
}

// Resolve implements ResolverT.Resolve() for void.
func (r S) Resolve(err error) {
	r.s.Resolve(nil, err)
}

// Forward implements ResolverT.Forward() for void.
func (r S) Forward(next R) {
	r.s.Forward(next.ResultT)
}
