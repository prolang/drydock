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

	"github.com/prolang/drydock/runtime/base/test"
)

// ReflectSuite is the test suite for reflection extensions in base.
type ReflectSuite struct {
	test.Suite
}

// TestReflectSuite runs the test suite for reflection extensions in base.
func TestReflectSuite(t *testing.T) {
	test.RunSuite(t, new(ReflectSuite))
}

// GetMethodName verifies the GetMethodName() function returns the correct string.
func (t *ReflectSuite) GetMethodName() {
	expected := "GetMethodName"
	if s := GetMethodName(); s != expected {
		t.Errorf("Invalid method name: Got: %s, Want: %s", s, expected)
	}
}
