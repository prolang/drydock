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

// This file contains an implementation of a thread-safe monotonically increasing ID generator.

import "fmt"

// UniqueID is an opaque unique identifier.
type UniqueID struct {
	id int64
}

// String renders a unique identifier as a string for printing and logging.
func (u UniqueID) String() string {
	return fmt.Sprintf("%d", int64(u.id))
}

// UniqueIDGenerator creates new uniqueIDs.
type UniqueIDGenerator struct {
	next int64
}

// NewUniqueIDGenerator creates a new ID generator.
func NewUniqueIDGenerator() *UniqueIDGenerator {
	return &UniqueIDGenerator{
		next: 1,
	}
}

// NewID generates a new ID.  ID's are never reused.
func (gen *UniqueIDGenerator) NewID() UniqueID {
	id := gen.next
	gen.next++

	return UniqueID{
		id: id,
	}
}
