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

// Package assert contains static methods for expression invariant conditions in code.
package assert

import log "github.com/golang/glog"

// True asserts an invariant must be true.  Panics if the condition is false.
func True(condition bool, format string, a ...interface{}) {
	if !condition {
		log.Fatalf("precondition failed: "+format, a...)
	}
}

// False asserts an invariant must be false.  Panics if the condition is true.
func False(condition bool, format string, a ...interface{}) {
	if condition {
		log.Fatalf("precondition failed: "+format, a...)
	}
}
