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

package base

import (
	"testing"

	"github.com/prolang/drydock/runtime/base/assert"
	"github.com/prolang/drydock/runtime/base/test"
)

// AssertSuite is the test suite for asserts.
type AssertSuite struct {
	test.Suite
}

// TestAssertSuite runs the assert test suite.
func TestAssertSuite(t *testing.T) {
	test.RunSuite(t, new(AssertSuite))
}

// Coverage provides minimal code coverage on the non-failing assert paths.
func (t *AssertSuite) Coverage() {
	assert.True(true, "Code coverage for assert.True()")
	assert.False(false, "Code coverage for assert.False()")
}
