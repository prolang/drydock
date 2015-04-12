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

// The types async.StringR and async.StringS represent a asynchronous computation with a string
// return value.  This represents an asynchronous computation of the form:
//
//   func() (string, error)
//
// See async.R (async_r.go) for a description of the model of computation.

import (
	"fmt"
	"reflect"
)

// StringR tracks the completion progress of an asynchronous computation.
type StringR struct {
	ResultT
}

// Type implements AwaitableT.Type().
func (StringR) Type() reflect.Type {
	return reflect.TypeOf((*string)(nil)).Elem()
}

// NewStringR allocates a new result.
func NewStringR(runner Runner) (StringR, StringS) {
	r, s := runner.NewResultT()
	return StringR{r}, StringS{s}
}

// NewStringError returns an unassociated already failed result.
func NewStringError(runner Runner, err error) StringR {
	r, s := NewStringR(runner)
	s.Fail(err)
	return r
}

// NewStringErrorf returns an unassociated already failed result.
func NewStringErrorf(runner Runner, format string, a ...interface{}) StringR {
	return NewStringError(runner, fmt.Errorf(format, a...))
}

// WhenString implements When.  See async.When().
func WhenString(r AwaitableT, f interface{}) StringR {
	reflectTypeStringR := reflect.TypeOf((*StringR)(nil)).Elem()
	reflectTypeString := reflect.TypeOf((*string)(nil)).Elem()
	return StringR{r.WhenT(r.Type(), reflectTypeString, reflectTypeStringR, f)}
}

// FinallyString implements Finally.  See async.Finally().
func FinallyString(r StringR, f func()) StringR {
	return WhenString(r, func(interface{}, error) StringR {
		f()
		return r
	})
}

// StringS is used to complete an asynchronous computation.
type StringS struct {
	s ResolverT
}

// Complete implements ResolverT.Complete().
func (r StringS) Complete(val string) {
	r.s.Complete(val)
}

// Fail implements ResolverT.Fail().
func (r StringS) Fail(err error) {
	r.s.Fail(err)
}

// Resolve implements ResolverT.Resolve().
func (r StringS) Resolve(val string, err error) {
	r.s.Resolve(val, err)
}

// Forward implements ResolverT.Forward().
func (r StringS) Forward(next StringR) {
	r.s.Forward(next.ResultT)
}
