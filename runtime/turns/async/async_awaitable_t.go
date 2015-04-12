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

// AwaitableT allows polymorphic treatment of ResultT subtypes that all allow their value to be
// waited for via a When expression.  All AwaitableT have a base type that corresponding to the
// concrete type of the value used to successfully complete the result.  Each ResultT subtype
// has a different base type.  E.g.
//
//   ResultT Type              Base Type
//     R                         "void"
//     StringR                   string
//
// ResultT subtypes are expected to embed a ResultT and then implement the Type() method.
type AwaitableT interface {
	// Base returns the same AwaitableT casted-up to untyped async result.
	Base() ResultT

	// Type returns the reflection type of the base type of the async result.
	Type() reflect.Type

	// WhenT schedules a function f that receives the outcome of a previous computation.  The function
	// f may optional take up to two input parameters:
	//
	//   value: contains the typed result of the previous computation if it was successful
	//   err: contains a non-nil error if the previous computation failed.
	//
	// If the previous computation failed value is undefined.
	//
	// If f does not take a value then the previous computation's result is discarded.
	//
	// If f does not take an error then f is NOT called in the event that the previous computation
	// failed and instead the result of the When expression is immediately failed with the same error.
	//
	// The function f's return value may take one of three forms:
	//
	//   concrete value
	//   concrete value and error
	//   async result
	//
	// If either of the first two forms then a resolved async result is created automatically from the
	// return values.  If an async result then the actual value becomes the result of the When
	// expression.
	//
	// The result of a When expression is always a typed async result of the same base type as f's
	// return value.  If f fails then the result of the When expression will eventually fail with the
	// same error.
	WhenT(in, out, outR reflect.Type, f interface{}) ResultT
}
